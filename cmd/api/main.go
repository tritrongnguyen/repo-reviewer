package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/tritrongnguyen/repo-reviewer.git/internal/server"
	"github.com/tritrongnguyen/repo-reviewer.git/pkg/logger"
	"go.uber.org/zap"
)

func gracefulShutdown(s *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logger.Log.Info("Shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Force shutdown", zap.Error(err))
	}

	logger.Log.Info("Server existing...")

	done <- true
}

func main() {
	_ = logger.Init(true)

	server := server.NewServer()

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err := server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done

	logger.Log.Info("Graceful shutdown complete.")
}
