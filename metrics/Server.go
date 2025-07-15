// Server.go

package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/util/httpwrapper"
)

type Server struct {
	httpServer *http.Server
	metricCfg  Metrics
	nfLogger   *logrus.Entry
}

// Initializes a new HTTP server instance and associate the prometheus handler to it
func NewServer(initMetrics InitMetrics, tlsKeyLogPath string, nfLogger *logrus.Entry) (*Server, error) {
	if nfLogger == nil {
		return nil, errors.New("nfLogger cannot be nil")
	}

	mux := http.NewServeMux()
	reg := Init(initMetrics)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	bindAddr := initMetrics.GetMetricsInfo().BindingIPv4
	nfLogger.Infof("Binding addr: [%s]", bindAddr)

	httpServer, err := httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, mux)
	if err != nil {
		nfLogger.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}

	s := &Server{
		httpServer: httpServer,
		metricCfg:  initMetrics.metricsInfo,
		nfLogger:   nfLogger,
	}

	return s, nil
}

// Configure the server to handle http requests
func (s *Server) ListenAndServe() {
	s.nfLogger.Infof("Starting HTTP server on %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.nfLogger.Errorf("Server error: %v", err)
	}
}

// Configure the server to handle https requests
func (s *Server) ListenAndServeTLS() {
	tlsKeyPath, tlsPemPath := s.metricCfg.Tls.Key, s.metricCfg.Tls.Pem

	err := s.httpServer.ListenAndServeTLS(tlsKeyPath, tlsPemPath)
	if err != nil {
		s.nfLogger.Errorf("Server error: %v", err)
	}
}

func (s *Server) startServer(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			s.nfLogger.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		wg.Done()
	}()

	var err error

	s.nfLogger.Infof("Start Metrics server (listen on %s)", s.httpServer.Addr)
	scheme := s.metricCfg.Scheme

	switch scheme {
	case "http":
		err = s.httpServer.ListenAndServe()
	case "https":
		tlsKeyPath, tlsPemPath := s.metricCfg.Tls.Key, s.metricCfg.Tls.Pem
		err = s.httpServer.ListenAndServeTLS(tlsKeyPath, tlsPemPath)
	default:
		err = fmt.Errorf("no support this scheme[%s]", scheme)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.nfLogger.Errorf("Metrics server error: %v", err)
	}
	s.nfLogger.Warnf("Metrics server (listen on %s) stopped", s.httpServer.Addr)
}

func (s *Server) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	go s.startServer(wg)
}

func (s *Server) Stop() {
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer != nil {
		s.nfLogger.Infof("Stop server (listen on %s)", s.httpServer.Addr)
		toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(toCtx); err != nil {
			s.nfLogger.Errorf("Could not close server: %#v", err)
		}
	}
}
