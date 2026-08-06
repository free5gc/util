package nfheartbeat

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

type fakeRegistrar struct {
	updateNf      models.NrfNfManagementNfProfile
	updatePd      *models.ProblemDetails
	updateErr     error
	updateDelay   time.Duration
	panicOnUpdate bool
	updateCalls   int

	registerTimer int32
	registerErr   error
	registerCalls int
}

func (f *fakeRegistrar) UpdateNFInstance(_ context.Context, _ []models.PatchItem) (
	models.NrfNfManagementNfProfile, *models.ProblemDetails, error,
) {
	f.updateCalls++
	if f.panicOnUpdate {
		panic("update panic")
	}
	if f.updateDelay > 0 {
		time.Sleep(f.updateDelay)
	}
	return f.updateNf, f.updatePd, f.updateErr
}

func (f *fakeRegistrar) RegisterNFInstance(_ context.Context) (int32, error) {
	f.registerCalls++
	return f.registerTimer, f.registerErr
}

func testLog() *logrus.Entry {
	log := logrus.New()
	log.SetOutput(io.Discard)
	return logrus.NewEntry(log)
}

func newTestRunner(t *testing.T, registrar Registrar, fallbackTimer func() int32) *Runner {
	t.Helper()

	runner, err := NewRunner(registrar, fallbackTimer, testLog())
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}
	return runner
}

// configFallback mirrors the config getter every NF passes: the configured
// value when set, DefaultTimer otherwise.
func configFallback(configTimer int32) func() int32 {
	return func() int32 {
		if configTimer > 0 {
			return configTimer
		}
		return DefaultTimer
	}
}

func TestNewRunnerValidation(t *testing.T) {
	if _, err := NewRunner(nil, configFallback(0), testLog()); err == nil {
		t.Error("NewRunner must reject a nil registrar")
	}
	if _, err := NewRunner(&fakeRegistrar{}, configFallback(0), nil); err == nil {
		t.Error("NewRunner must reject a nil log")
	}
	if _, err := NewRunner(&fakeRegistrar{}, nil, testLog()); err != nil {
		t.Errorf("NewRunner with a nil fallbackTimer: %v", err)
	}
}

func TestInterval(t *testing.T) {
	tests := []struct {
		name     string
		nrfTimer int32
		fallback func() int32
		want     time.Duration
	}{
		{
			name:     "no timer at all falls back to the default",
			fallback: configFallback(0),
			want:     time.Duration(DefaultTimer) * time.Second,
		},
		{
			name:     "negative NRF timer falls back to the default",
			nrfTimer: -5,
			fallback: configFallback(0),
			want:     time.Duration(DefaultTimer) * time.Second,
		},
		{
			name:     "no NRF timer falls back to the configured one",
			fallback: configFallback(45),
			want:     45 * time.Second,
		},
		{
			name:     "NRF timer takes precedence over the configured one",
			nrfTimer: 10,
			fallback: configFallback(45),
			want:     10 * time.Second,
		},
		{
			name:     "timer at the NRF profile cap is adopted",
			nrfTimer: MaxTimer,
			fallback: configFallback(0),
			want:     time.Duration(MaxTimer) * time.Second,
		},
		{
			name:     "zero fallback clamps to the default",
			fallback: func() int32 { return 0 },
			want:     time.Duration(DefaultTimer) * time.Second,
		},
		{
			name:     "nil fallback clamps to the default",
			fallback: nil,
			want:     time.Duration(DefaultTimer) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := newTestRunner(t, &fakeRegistrar{}, tt.fallback)
			runner.timer = tt.nrfTimer

			if got := runner.interval(); got != tt.want {
				t.Errorf("interval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTick(t *testing.T) {
	// The NRF assigned currentTimer at registration time and answers every
	// re-registration with reregisterTimer.
	const (
		currentTimer    = 10
		reregisterTimer = 30
	)

	notFoundErr := openapi.GenericOpenAPIError{ErrorStatus: http.StatusNotFound}
	serverErr := openapi.GenericOpenAPIError{ErrorStatus: http.StatusInternalServerError}

	tests := []struct {
		name          string
		updateNf      models.NrfNfManagementNfProfile
		updatePd      *models.ProblemDetails
		updateErr     error
		panicOnUpdate bool
		ticks         int
		wantInterval  time.Duration
		wantFailures  int
		wantRegister  bool
	}{
		{
			name:         "200 adopts the returned timer",
			updateNf:     models.NrfNfManagementNfProfile{HeartBeatTimer: 25},
			ticks:        1,
			wantInterval: 25 * time.Second,
		},
		{
			name:         "204 keeps the current timer",
			ticks:        1,
			wantInterval: currentTimer * time.Second,
		},
		{
			name:         "200 above the NRF profile cap is capped",
			updateNf:     models.NrfNfManagementNfProfile{HeartBeatTimer: MaxTimer + 1},
			ticks:        1,
			wantInterval: time.Duration(MaxTimer) * time.Second,
		},
		{
			name:         "404 re-registers right away",
			updatePd:     &models.ProblemDetails{Status: http.StatusNotFound},
			updateErr:    notFoundErr,
			ticks:        1,
			wantInterval: reregisterTimer * time.Second,
			wantRegister: true,
		},
		{
			name:         "404 as problem details without an error re-registers",
			updatePd:     &models.ProblemDetails{Status: http.StatusNotFound},
			ticks:        1,
			wantInterval: reregisterTimer * time.Second,
			wantRegister: true,
		},
		{
			name:         "404 as a pointer error re-registers",
			updateErr:    &openapi.GenericOpenAPIError{ErrorStatus: http.StatusNotFound},
			ticks:        1,
			wantInterval: reregisterTimer * time.Second,
			wantRegister: true,
		},
		{
			name:         "failures below the threshold only accumulate",
			updatePd:     &models.ProblemDetails{Status: http.StatusInternalServerError},
			updateErr:    serverErr,
			ticks:        reregisterFailureThreshold - 1,
			wantInterval: currentTimer * time.Second,
			wantFailures: reregisterFailureThreshold - 1,
		},
		{
			name:         "threshold re-registers and resets the failures",
			updatePd:     &models.ProblemDetails{Status: http.StatusInternalServerError},
			updateErr:    serverErr,
			ticks:        reregisterFailureThreshold,
			wantInterval: reregisterTimer * time.Second,
			wantRegister: true,
		},
		{
			name:          "panicking attempts count toward the threshold",
			panicOnUpdate: true,
			ticks:         reregisterFailureThreshold,
			wantInterval:  reregisterTimer * time.Second,
			wantRegister:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrar := &fakeRegistrar{
				updateNf:      tt.updateNf,
				updatePd:      tt.updatePd,
				updateErr:     tt.updateErr,
				panicOnUpdate: tt.panicOnUpdate,
				registerTimer: reregisterTimer,
			}
			runner := newTestRunner(t, registrar, configFallback(0))
			runner.timer = currentTimer

			for range tt.ticks {
				runner.tick(t.Context())
			}

			if got := runner.interval(); got != tt.wantInterval {
				t.Errorf("interval() = %v, want %v", got, tt.wantInterval)
			}
			if runner.failures != tt.wantFailures {
				t.Errorf("failures = %d, want %d", runner.failures, tt.wantFailures)
			}
			wantRegisterCalls := 0
			if tt.wantRegister {
				wantRegisterCalls = 1
			}
			if registrar.registerCalls != wantRegisterCalls {
				t.Errorf("registerCalls = %d, want %d", registrar.registerCalls, wantRegisterCalls)
			}
			if registrar.updateCalls != tt.ticks {
				t.Errorf("updateCalls = %d, want %d", registrar.updateCalls, tt.ticks)
			}
		})
	}
}

func TestSuccessResetsFailures(t *testing.T) {
	registrar := &fakeRegistrar{
		updatePd: &models.ProblemDetails{Status: http.StatusInternalServerError},
	}
	runner := newTestRunner(t, registrar, configFallback(0))

	runner.tick(t.Context())
	runner.tick(t.Context())
	registrar.updatePd = nil
	runner.tick(t.Context())

	if runner.failures != 0 {
		t.Errorf("failures = %d, want 0 after a successful heartbeat", runner.failures)
	}
	if registrar.registerCalls != 0 {
		t.Errorf("registerCalls = %d, want 0: the threshold was never reached", registrar.registerCalls)
	}
}

func TestTickSurvivesPanic(t *testing.T) {
	registrar := &fakeRegistrar{panicOnUpdate: true}
	runner := newTestRunner(t, registrar, configFallback(0))

	runner.tick(t.Context())

	if runner.failures != 1 {
		t.Errorf("failures = %d, want 1 after a recovered panic", runner.failures)
	}
}

func TestRecoverRegistrationSkippedOnShutdown(t *testing.T) {
	registrar := &fakeRegistrar{}
	runner := newTestRunner(t, registrar, configFallback(0))

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if runner.recoverRegistration(ctx) {
		t.Error("recoverRegistration must report a failure when shutting down")
	}
	if registrar.registerCalls != 0 {
		t.Error("no registration must be sent when shutting down")
	}
}

func TestTickSkippedOnShutdown(t *testing.T) {
	registrar := &fakeRegistrar{}
	runner := newTestRunner(t, registrar, configFallback(0))

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	runner.tick(ctx)

	if registrar.updateCalls != 0 {
		t.Error("no heartbeat must be sent when shutting down")
	}
	if runner.failures != 0 {
		t.Errorf("failures = %d, want 0 for a skipped tick", runner.failures)
	}
}

func TestLifecycle(t *testing.T) {
	t.Run("wait returns when the heartbeat never started", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			waitStopped(t, newTestRunner(t, &fakeRegistrar{}, configFallback(0)))
		})
	})

	t.Run("heartbeat stops when the context is cancelled", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runner := newTestRunner(t, &fakeRegistrar{}, configFallback(0))
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(t.Context())
			runner.Start(ctx, &wg, 0)
			cancel()

			waitStopped(t, runner)
			wg.Wait()
		})
	})

	t.Run("start with a cancelled context does nothing", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			registrar := &fakeRegistrar{}
			runner := newTestRunner(t, registrar, configFallback(0))
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			runner.Start(ctx, &wg, 0)

			waitStopped(t, runner)
			wg.Wait()
			if registrar.updateCalls != 0 {
				t.Errorf("updateCalls = %d, want 0", registrar.updateCalls)
			}
		})
	})
}

func TestSecondStartIsIgnored(t *testing.T) {
	const nrfTimer = 1

	synctest.Test(t, func(t *testing.T) {
		registrar := &fakeRegistrar{}
		runner := newTestRunner(t, registrar, configFallback(0))
		var wg sync.WaitGroup

		ctx, cancel := context.WithCancel(t.Context())
		runner.Start(ctx, &wg, nrfTimer)
		runner.Start(ctx, &wg, nrfTimer)

		time.Sleep(nrfTimer * time.Second)
		synctest.Wait()

		cancel()
		waitStopped(t, runner)
		wg.Wait()

		if registrar.updateCalls != 1 {
			t.Errorf("updateCalls = %d, want 1: a second Start must not add a loop", registrar.updateCalls)
		}
	})
}

func TestZeroFallbackStillHeartbeats(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		registrar := &fakeRegistrar{}
		runner := newTestRunner(t, registrar, func() int32 { return 0 })
		var wg sync.WaitGroup

		ctx, cancel := context.WithCancel(t.Context())
		runner.Start(ctx, &wg, 0)

		// Without the DefaultTimer clamp this would panic in NewTicker(0) and
		// the heartbeat would silently never start.
		time.Sleep(time.Duration(DefaultTimer) * time.Second)
		synctest.Wait()

		cancel()
		waitStopped(t, runner)
		wg.Wait()

		if registrar.updateCalls != 1 {
			t.Errorf("updateCalls = %d, want 1 tick at the default interval", registrar.updateCalls)
		}
	})
}

func TestLoopAdoptsNewInterval(t *testing.T) {
	// The NRF assigned initialTimer at registration time and answers the first
	// heartbeat with adoptedTimer.
	const (
		initialTimer = 1
		adoptedTimer = 2
	)

	synctest.Test(t, func(t *testing.T) {
		registrar := &fakeRegistrar{
			updateNf: models.NrfNfManagementNfProfile{HeartBeatTimer: adoptedTimer},
		}
		runner := newTestRunner(t, registrar, configFallback(0))

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(t.Context())
		runner.Start(ctx, &wg, initialTimer)

		// Fake clock: this returns as soon as the loop has served the tick.
		time.Sleep(initialTimer * time.Second)
		synctest.Wait()

		// Read the fake and the adopted interval only once the goroutine that
		// owns them has exited.
		cancel()
		waitStopped(t, runner)
		wg.Wait()

		if registrar.updateCalls == 0 {
			t.Fatal("the heartbeat loop sent no PATCH")
		}
		if got := runner.interval(); got != adoptedTimer*time.Second {
			t.Errorf("interval() = %v, want %ds adopted by the loop", got, adoptedTimer)
		}
	})
}

// TestWaitCoversInFlightTick proves the guarantee deregistration relies on:
// Wait may not return while a heartbeat PATCH is still in flight, even after
// the shutdown has been signalled.
func TestWaitCoversInFlightTick(t *testing.T) {
	const (
		interval    = 1
		updateDelay = 3 * time.Second
	)

	synctest.Test(t, func(t *testing.T) {
		registrar := &fakeRegistrar{updateDelay: updateDelay}
		runner := newTestRunner(t, registrar, configFallback(0))
		var wg sync.WaitGroup

		ctx, cancel := context.WithCancel(t.Context())
		runner.Start(ctx, &wg, interval)

		// Enter the tick, then shut down while the PATCH is still in flight.
		time.Sleep(interval * time.Second)
		synctest.Wait()
		cancel()

		stopped := make(chan struct{})
		go func() {
			runner.Wait()
			close(stopped)
		}()
		synctest.Wait()
		select {
		case <-stopped:
			t.Fatal("Wait returned while a heartbeat was still in flight")
		default:
		}

		// The PATCH completes, the loop sees the cancelled ctx and exits.
		time.Sleep(updateDelay)
		synctest.Wait()
		select {
		case <-stopped:
		default:
			t.Fatal("Wait did not return after the in-flight heartbeat completed")
		}
		wg.Wait()
	})
}

// waitStopped fails the test when the heartbeat goroutine outlives the wait
// that deregistration relies on. It must run inside a synctest bubble: the
// deadlock it guards against shows up as a blocked bubble, not as a timeout.
func waitStopped(t *testing.T, runner *Runner) {
	t.Helper()

	stopped := make(chan struct{})
	go func() {
		runner.Wait()
		close(stopped)
	}()

	synctest.Wait()
	select {
	case <-stopped:
	default:
		t.Fatal("Wait did not return")
	}
}
