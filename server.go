package k8server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(ctx context.Context, public *http.ServeMux, options ...Option) error {
	config := defaults()

	for _, option := range options {
		option.apply(config)
	}

	publicSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: public,
	}

	errs := make(chan error, 1)
	go func(err chan<- error) {
		err <- publicSrv.ListenAndServe()
	}(errs)

	mgmt := http.NewServeMux()
	mgmt.Handle("GET /metrics", promhttp.Handler())
	mgmt.Handle("GET /livez", Livez())
	mgmt.Handle("GET /readyz", Readyz(&publicSrv))
	mgmtSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.ManagementPort),
		Handler: mgmt,
	}

	go func(err chan<- error) {
		err <- mgmtSrv.ListenAndServe()
	}(errs)

	slog.InfoContext(ctx, "serving public endpoints", "uri", fmt.Sprintf("http://localhost:%d", config.Port))
	slog.InfoContext(ctx, "serving metrics endpoint", "uri", fmt.Sprintf("http://localhost:%d/metrics", config.ManagementPort))
	slog.InfoContext(ctx, "serving liveness endpoint", "uri", fmt.Sprintf("http://localhost:%d/readyz", config.ManagementPort))
	slog.InfoContext(ctx, "serving readiness endpoint", "uri", fmt.Sprintf("http://localhost:%d/livez", config.ManagementPort))

	// Wait for a fatal server error, an interrupt signal, or a termination signal.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-signals:
		slog.InfoContext(ctx, "received signal", "signal", sig)

		slog.InfoContext(ctx, "waiting for remaining requests to finish", "timeout", config.Timeout)
		terminationCtx, cancel := context.WithTimeout(ctx, config.Timeout)
		defer cancel()

		if err := publicSrv.Shutdown(terminationCtx); err != nil {
			slog.ErrorContext(terminationCtx, "server shutdown error", "error", err)
			return err
		}

		slog.InfoContext(ctx, "graceful shutdown complete")
	case err := <-errs:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "server error", "error", err)
			return err
		}
	}

	return nil
}
