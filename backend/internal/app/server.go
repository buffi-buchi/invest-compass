package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api/middleware"
	authapi "github.com/buffi-buchi/invest-compass/backend/internal/api/v1/auth"
	portfolioapi "github.com/buffi-buchi/invest-compass/backend/internal/api/v1/portfolio"
	userapi "github.com/buffi-buchi/invest-compass/backend/internal/api/v1/user"
	"github.com/buffi-buchi/invest-compass/backend/internal/domain/auth"
	"github.com/buffi-buchi/invest-compass/backend/internal/domain/user"
	"github.com/buffi-buchi/invest-compass/backend/internal/provider/jwt"
	"github.com/buffi-buchi/invest-compass/backend/internal/provider/postgres"
)

func RunServer() error {
	// Configure logger.
	logger, err := NewLogger()
	if err != nil {
		return err
	}

	defer logger.Sync()

	// Read configuration.
	const envVarConfigPath = "CONFIG_PATH"

	configPath := os.Getenv(envVarConfigPath)
	if configPath == "" {
		logger.Error("Config path is required")
		return errors.New("config path is required")
	}

	config, err := ReadConfig(configPath)
	if err != nil {
		logger.Error("Failed to read config", zap.Error(err))
		return fmt.Errorf("read config: %w", err)
	}

	// Configure database connection.
	dbConfig, err := pgxpool.ParseConfig(config.Postgres.GetConnectionString())
	if err != nil {
		logger.Error("Failed to parse database config", zap.Error(err))
		return fmt.Errorf("parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), dbConfig)
	if err != nil {
		logger.Error("Failed to create database pool", zap.Error(err))
		return fmt.Errorf("create pool: %w", err)
	}

	// Configure providers.
	jwtProvider := jwt.NewProvider([]byte("secretkey"), "server", 30*time.Minute)

	// Configure stores.
	userStore := postgres.NewUserStore(pool)
	portfolioStore := postgres.NewPortfolioStore(pool)

	// Configure services.
	authService := auth.NewService(userStore, jwtProvider)
	userService := user.NewService(userStore)

	// Configure middlewares.
	authMiddleware := middleware.NewAuthMiddleware(jwtProvider)

	// Configure controllers.
	authController := authapi.NewImplementation(authService, logger)
	userController := userapi.NewImplementation(userService, logger)
	portfolioController := portfolioapi.NewImplementation(portfolioStore, authMiddleware, logger)

	// Start HTTP servers.

	mux := chi.NewMux()

	authController.Register(mux)
	userController.Register(mux)
	portfolioController.Register(mux)

	server := http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: mux,
	}

	mux = chi.NewMux()
	mux.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	debugServer := http.Server{
		Addr:    ":" + config.DebugServer.Port,
		Handler: mux,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		logger.Info("Starting server", zap.String("address", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start server", zap.Error(err))
			cancel()
		}
	}()

	go func() {
		logger.Info("Starting debug server", zap.String("address", debugServer.Addr))

		if err := debugServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start debug server", zap.Error(err))
			cancel()
		}
	}()

	// Graceful shutdown.
	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server")
	_ = server.Shutdown(ctx)

	logger.Info("Shutting down debug server")
	_ = debugServer.Shutdown(ctx)

	logger.Info("Application stopped")

	return nil
}
