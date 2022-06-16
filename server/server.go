package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	notifySignal   = signal.Notify
	serverShutdown = func(server *http.Server, ctx context.Context) error {
		return server.Shutdown(ctx)
	}
)

// Start starts the http server.
func Start(handler http.Handler, addr string) error {
	server, shutdown := startServer(handler, addr)
	shutdownOnInterruptSignal(server, 2*time.Second, shutdown)
	return waitForServerToClose(shutdown)
}

func startServer(handler http.Handler, addr string) (*http.Server, chan error) {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	shutdown := make(chan error)
	go func() {
		err := srv.ListenAndServe()
		shutdown <- err
	}()
	return srv, shutdown
}

func shutdownOnInterruptSignal(server *http.Server, timeout time.Duration, shutdown chan<- error) {
	interrupt := make(chan os.Signal, 1)
	notifySignal(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := serverShutdown(server, ctx); err != nil {
			shutdown <- err
		}
	}()
}

func waitForServerToClose(shutdown <-chan error) error {
	err := <-shutdown
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
