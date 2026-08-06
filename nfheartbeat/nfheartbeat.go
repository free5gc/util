// Package nfheartbeat keeps an NF profile registered with the NRF by sending
// the periodic NF heartbeat defined in 3GPP TS 29.510 clause 5.2.2.3.2.
package nfheartbeat

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

const (
	// DefaultTimer is the fallback heartbeat interval in seconds, applied when
	// neither the NRF nor fallbackTimer supplies a positive value. Keep it
	// equal to the default interval of the free5gc NRF heartbeat enforcement
	// (nrf branch feat/nf-heartbeat): a longer fallback would cross the
	// deadline past which that NRF suspends a silent instance.
	DefaultTimer int32 = 10

	// MinTimer and MaxTimer mirror the heartBeatTimer bounds of the free5gc
	// NRF profile validation (nrf branch feat/nf-heartbeat). NF config
	// validators should stay within them.
	MinTimer int32 = 1
	MaxTimer int32 = 3600
)

// hbStatus is the outcome of one heartbeat attempt.
type hbStatus int

const (
	hbStatusOk hbStatus = iota
	hbStatusNotFound
	hbStatusFailed
)

// reregisterFailureThreshold is the consecutive heartbeat failure count that
// triggers a full re-registration. It covers NRF database loss behind OAuth2,
// where the token request fails before the PATCH can observe a 404.
const reregisterFailureThreshold = 3

// PatchItems returns the NF heartbeat body from 3GPP TS 29.510 clause 5.2.2.3.2.
func PatchItems() []models.PatchItem {
	return []models.PatchItem{{
		Op:    models.PatchOperation_REPLACE,
		Path:  "/nfStatus",
		Value: models.NrfNfManagementNfStatus_REGISTERED,
	}}
}

// Registrar executes the NRF requests the Runner decides to send: the Runner
// picks when to heartbeat or re-register, the Registrar carries it out. Each
// NF implements it over its own consumer, so the SBI clients, the OAuth2
// token handling and the NF profile stay on the NF side.
type Registrar interface {
	// UpdateNFInstance sends the heartbeat PATCH. A non-2xx answer must come
	// back as a non-nil ProblemDetails or as an error carrying the raw
	// openapi.GenericOpenAPIError, so the Runner can classify a 404. ctx is
	// cancelled on shutdown only, so implementations should still bound the
	// request to less than the heartbeat interval.
	UpdateNFInstance(ctx context.Context, patchItems []models.PatchItem) (
		models.NrfNfManagementNfProfile, *models.ProblemDetails, error)
	// RegisterNFInstance re-registers the NF profile with the NRF and returns
	// the heartBeatTimer it assigned, in seconds, 0 when it assigned none.
	// It may retry internally until it succeeds or ctx is done.
	RegisterNFInstance(ctx context.Context) (int32, error)
}

// Runner sends the periodic NF heartbeat and re-registers the NF profile when
// the NRF has lost it. Only the first Start call has any effect.
type Runner struct {
	registrar     Registrar
	fallbackTimer func() int32
	log           *logrus.Entry

	started atomic.Bool
	wg      sync.WaitGroup

	// Seeded by Start, then owned by the heartbeat goroutine.
	// timer is the interval in seconds last assigned by the NRF.
	timer    int32
	failures int
}

// NewRunner returns a Runner driving registrar. fallbackTimer supplies the
// interval in seconds while the NRF has not assigned one, typically the NF's
// config getter; a nil fallbackTimer or a non-positive value falls back to
// DefaultTimer.
func NewRunner(registrar Registrar, fallbackTimer func() int32, log *logrus.Entry) (*Runner, error) {
	if registrar == nil {
		return nil, errors.New("registrar cannot be nil")
	}
	if log == nil {
		return nil, errors.New("log cannot be nil")
	}
	return &Runner{
		registrar:     registrar,
		fallbackTimer: fallbackTimer,
		log:           log,
	}, nil
}

// Start launches the periodic NF heartbeat toward the NRF. It must be called
// after a successful NF registration; nrfTimer is the heartBeatTimer assigned
// there, in seconds, 0 when the NRF assigned none. wg tracks the goroutine
// app-wide, while Wait covers just the heartbeat. Calls after the first, or
// with ctx already cancelled, do nothing.
func (r *Runner) Start(ctx context.Context, wg *sync.WaitGroup, nrfTimer int32) {
	// A Start that raced shutdown must not begin heartbeating.
	if ctx.Err() != nil {
		return
	}
	if !r.started.CompareAndSwap(false, true) {
		r.log.Warnln("NF heartbeat already started, ignoring Start")
		return
	}
	r.timer = r.capTimer(nrfTimer)
	// The external Add comes first: a nil wg then panics before the private
	// WaitGroup is touched, so Wait cannot block on a goroutine never started.
	wg.Add(1)
	r.wg.Add(1)
	go func() {
		defer wg.Done()
		defer r.wg.Done()
		r.loop(ctx)
	}()
}

// Wait blocks until the heartbeat goroutine has exited, so that no heartbeat
// PATCH or re-registration PUT can reach the NRF after deregistration. It
// returns immediately when the heartbeat was never started.
func (r *Runner) Wait() {
	r.wg.Wait()
}

// loop sends the NF heartbeat at the current interval and re-adopts the
// heartBeatTimer carried by each answer, per 3GPP TS 29.510 clause 5.2.2.3.2.
func (r *Runner) loop(ctx context.Context) {
	// Ticks recover on their own; this one is the last resort for the loop
	// machinery itself. A panic escaping this goroutine would crash the whole
	// process. The heartbeat stays down after this recover fires: the NF keeps
	// serving and the NRF suspension makes the outage visible.
	defer func() {
		if p := recover(); p != nil {
			r.log.Errorf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	interval := r.interval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	r.log.Infof("NF heartbeat started, interval %v", interval)

	for {
		select {
		case <-ctx.Done():
			r.log.Infoln("NF heartbeat stopped")
			return
		case <-ticker.C:
			r.tick(ctx)
			if next := r.interval(); next != interval {
				interval = next
				ticker.Reset(interval)
				r.log.Infof("NF heartbeat interval updated to %v", interval)
			}
		}
	}
}

// tick sends one heartbeat and re-registers the NF profile when the NRF no
// longer knows it, or after reregisterFailureThreshold consecutive failures,
// panicking attempts included.
func (r *Runner) tick(ctx context.Context) {
	// Contains a panic escaping the re-registration path; the heartbeat
	// attempt recovers separately in guardedOnce.
	defer func() {
		if p := recover(); p != nil {
			r.log.Errorf("panic during NF re-registration: %v\n%s", p, string(debug.Stack()))
		}
	}()

	// A tick that fired alongside the shutdown must not reach the NRF.
	if ctx.Err() != nil {
		return
	}

	switch r.guardedOnce(ctx) {
	case hbStatusOk:
		r.failures = 0
		return
	case hbStatusNotFound:
		r.log.Warnln("NF profile not found on NRF, re-registering")
	case hbStatusFailed:
		r.failures++
		if r.failures < reregisterFailureThreshold {
			return
		}
		r.log.Warnf("%d consecutive NF heartbeat failures, re-registering", r.failures)
	}
	if r.recoverRegistration(ctx) {
		r.failures = 0
	}
}

// guardedOnce sends one heartbeat, turning a panic into a failure so that it
// costs one heartbeat instead of silently ending them all.
func (r *Runner) guardedOnce(ctx context.Context) (status hbStatus) {
	defer func() {
		if p := recover(); p != nil {
			r.log.Errorf("panic during NF heartbeat: %v\n%s", p, string(debug.Stack()))
			status = hbStatusFailed
		}
	}()
	return r.once(ctx)
}

// once sends one heartbeat PATCH: 200 carries the profile, so its
// heartBeatTimer is adopted, 204 carries nothing and the current interval
// stands, 404 means the NRF no longer holds the profile. A 200 whose
// heartBeatTimer is 0 also keeps the current interval: the int32 model cannot
// tell an explicit 0 from an absent field, so 0 must not read as disable the
// way the legacy openapi nrf/service.go helper reads it.
func (r *Runner) once(ctx context.Context) hbStatus {
	nf, problemDetails, err := r.registrar.UpdateNFInstance(ctx, PatchItems())
	if err == nil && problemDetails == nil {
		if nf.HeartBeatTimer > 0 {
			r.timer = r.capTimer(nf.HeartBeatTimer)
		}
		return hbStatusOk
	}
	if isNotFound(problemDetails, err) {
		return hbStatusNotFound
	}
	r.log.Warnf("NF heartbeat failed: pd=%+v err=%+v", problemDetails, err)
	return hbStatusFailed
}

// isNotFound reports whether the NRF answered with a 404, in any of the shapes
// a Registrar may deliver it.
func isNotFound(pd *models.ProblemDetails, err error) bool {
	if pd != nil && pd.Status == http.StatusNotFound {
		return true
	}
	var apiErr openapi.GenericOpenAPIError
	if errors.As(err, &apiErr) && apiErr.ErrorStatus == http.StatusNotFound {
		return true
	}
	var apiErrPtr *openapi.GenericOpenAPIError
	if errors.As(err, &apiErrPtr) && apiErrPtr.ErrorStatus == http.StatusNotFound {
		return true
	}
	return false
}

// recoverRegistration re-registers the NF profile with the NRF and reports
// whether it succeeded. It gives up while the NF is shutting down: a PUT after
// the deregistration would resurrect the profile.
func (r *Runner) recoverRegistration(ctx context.Context) bool {
	if ctx.Err() != nil {
		return false
	}
	timer, err := r.registrar.RegisterNFInstance(ctx)
	if err != nil {
		r.log.Errorf("NF re-registration aborted: %+v", err)
		return false
	}
	r.timer = r.capTimer(timer)
	r.log.Infoln("NF re-registered to NRF")
	return true
}

// capTimer caps an NRF-assigned interval at MaxTimer so a misbehaving NRF
// cannot park the heartbeat for hours. MinTimer needs no clamp: any positive
// int32 already satisfies it, and non-positive values keep their fallback
// meaning.
func (r *Runner) capTimer(timer int32) int32 {
	if timer > MaxTimer {
		r.log.Warnf("NRF-assigned heartbeat timer %ds capped to %ds", timer, MaxTimer)
		return MaxTimer
	}
	return timer
}

// interval is the interval assigned by the NRF, in seconds, falling back to
// fallbackTimer and then to DefaultTimer. It can never be non-positive: the
// ticker panics on zero.
func (r *Runner) interval() time.Duration {
	timer := r.timer
	if timer <= 0 && r.fallbackTimer != nil {
		timer = r.fallbackTimer()
	}
	if timer <= 0 {
		timer = DefaultTimer
	}
	return time.Duration(timer) * time.Second
}
