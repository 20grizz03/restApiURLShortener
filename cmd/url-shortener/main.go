package main

import (
	"github.com/20grizz03/restApiURLShortener/internal/config"
	"github.com/20grizz03/restApiURLShortener/internal/db/sqlite"
	"github.com/20grizz03/restApiURLShortener/internal/http-server/handlers/deleteURL"
	"github.com/20grizz03/restApiURLShortener/internal/http-server/handlers/redirect"
	save "github.com/20grizz03/restApiURLShortener/internal/http-server/handlers/url"
	mwLogger "github.com/20grizz03/restApiURLShortener/internal/http-server/middleware/logger"
	"github.com/20grizz03/restApiURLShortener/internal/lib/logger/hendlers/slogpretty"
	"github.com/20grizz03/restApiURLShortener/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("app started")
	log.Debug("app started in debug mode")

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init db", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	// middlewares

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, db))
		r.Delete("/{alias}", deleteURL.New(log, db))
	})

	router.Get("/{alias}", redirect.New(log, db))

	log.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("server error", sl.Err(err))
	}

	log.Error("server stopped", sl.Err(err))

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
