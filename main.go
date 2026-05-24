package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"mthan-go-starter/app/config"
	"mthan-go-starter/app/mods"
	"mthan-go-starter/app/routes/api"
	"mthan-go-starter/app/routes/api/hello"
	"mthan-go-starter/app/routes/post/action"
	"mthan-go-starter/app/services"
)

func main() {
	// 1. Load application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Logger Service
	logger := services.NewLoggerService(cfg.Env)
	logger.Info("Configuration loaded successfully")

	// 3. Initialize layers
	svc := services.NewProcessService()
	mux := http.NewServeMux()
	
	// Register individual route packages
	api.Register(mux, svc, logger)
	hello.Register(mux, svc, logger)
	action.Register(mux, svc, logger)
	
	logger.Info("Admin Panel is not compiled in this build")

	// Serve React client if enabled, otherwise register fallback landing page
	if cfg.Client {
		logger.Info("Serving client static files", "path", cfg.ClientBuildPath)
		fs := http.Dir(cfg.ClientBuildPath)
		fileServer := http.FileServer(fs)

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			f, err := fs.Open(r.URL.Path)
			if err != nil {
				// File does not exist, serve index.html for SPA routing fallback
				http.ServeFile(w, r, filepath.Join(cfg.ClientBuildPath, "index.html"))
				return
			}
			f.Close()
			fileServer.ServeHTTP(w, r)
		})
	} else {
		mods.Register(mux, svc, logger)
	}

	// 4. Setup HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. Handle Graceful Shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	go func() {
		<-sigChan
		logger.Info("Shutting down server gracefully...")

		// Shutdown context with a 10 second timeout limit
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Error("Graceful shutdown timed out. Forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP server Shutdown failed", "error", err)
			os.Exit(1)
		}
		serverStopCtx()
	}()

	// 6. Start HTTP Server
	logger.Info("Server is running", "addr", server.Addr, "mode", cfg.Env)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("ListenAndServe failed to start", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown goroutine to finish cleaning up
	<-serverCtx.Done()
	logger.Info("Server stopped cleanly")
}
