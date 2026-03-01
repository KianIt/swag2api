//go:generate go run ../cmd/swag2api -handler=s2aHandler
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title swag2api example documntation
// @version 1.0

var s2aHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("Server API not generated"))
})

func main() {
	server := createServer(":8080")

	if err := runServer(context.Background(), server, 10*time.Second); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func createServer(port string) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: loggingMiddleware(s2aHandler),
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			body []byte
			err  error
		)
		if body, err = io.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		log.Printf("HTTP Request from: %s, to: %s, method: %s, body: %s", r.RemoteAddr, r.URL.Path, r.Method, string(body))

		next.ServeHTTP(w, r)
	})
}

func runServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	errChan := make(chan error, 1)

	go func() {
		log.Printf("Server listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("listening error: %w", err)
	case <-stop:
		log.Println("Server shutdown signal received")
	case <-ctx.Done():
		log.Println("Context canceled, shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return fmt.Errorf("closing error: %w", closeErr)
		}
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Println("Server shutdown gracefully")

	return nil
}
