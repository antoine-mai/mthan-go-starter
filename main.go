package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mthan-go-starter/app/mods/admin"
	"mthan-go-starter/app/mods/client"
	routes_api_hello "mthan-go-starter/app/routes/api/hello"
	routes_api "mthan-go-starter/app/routes/api"
	routes_pages_contact "mthan-go-starter/app/routes/pages/contact"
	routes_pages "mthan-go-starter/app/routes/pages"
	routes_post_action "mthan-go-starter/app/routes/post/action"
	"mthan-go-starter/app/services"
)

func main() {
	// 1. Load application configuration
	cfg, err := services.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Logger Service
	logger := services.NewLoggerService(cfg.Env)
	logger.Info("Configuration loaded successfully")

	// 3. Initialize Database Service
	dbSvc := services.NewDatabaseService(cfg, logger)
	if _, err := dbSvc.Connect(); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbSvc.Close()

	// 4. Initialize layers
	svc := services.NewProcessService()
	mux := http.NewServeMux()
	
	// Register individual route packages dynamically
	routes_api_hello.Register(mux, svc, logger, cfg)
	routes_api.Register(mux, svc, logger, cfg)
	routes_post_action.Register(mux, svc, logger, cfg)
	if cfg.AdminActive && cfg.AdminPath != "" && cfg.AdminUser != "" && cfg.AdminPass != "" {
		admin.Register(mux, svc, logger, cfg)
	} else {
		logger.Info("Admin Panel is disabled (admin.active is false or missing credentials in config)")
	}
	// Serve React client if enabled, otherwise register dynamic pages, templates, or fallback landing page
	if cfg.Client {
		client.Register(mux, svc, logger, cfg)
	} else {
		// Register dynamic SSR HTML Page routes
		routes_pages_contact.Register(mux, svc, logger, cfg)
		routes_pages.Register(mux, svc, logger, cfg)
	}

	// 5. Setup HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 6. Handle Graceful Shutdown
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

	// 7. Start HTTP Server
	logger.Info("Server is running", "addr", server.Addr, "mode", cfg.Env)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("ListenAndServe failed to start", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown goroutine to finish cleaning up
	<-serverCtx.Done()
	logger.Info("Server stopped cleanly")
}
