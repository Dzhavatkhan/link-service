package main

import (
	"context"
	"net/http"

	"link-service/internal/config"
	"link-service/internal/http-server/handlers/delete"
	"link-service/internal/http-server/handlers/redirect"
	"link-service/internal/http-server/handlers/update"
	"link-service/internal/http-server/handlers/url/save"
	"link-service/internal/lib/logger/handlers/slogpretty"
	"link-service/internal/lib/logger/logslog"
	postgres "link-service/internal/storage"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("started link-service", slog.String("env", cfg.Env));
	storage, err := postgres.New(ctx);
	if err != nil {
		log.Error("failed to connection db", logslog.Error(err));
		os.Exit(1);
	}
	//TODO: init router chi
	router := chi.NewRouter();
	//TODO: init middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use( middleware.Recoverer);
	router.Use(middleware.URLFormat);

	router.Post("/url", save.New(log, storage));
	router.Get("/{alias}", redirect.New(log, storage));
	router.Patch("/update/{id}", update.New(log, storage));	
	router.Delete("/delete/{id}", delete.New(log, storage));

	
	log.Info("starting server", slog.String("address", cfg.Address))
	
	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err);
	}

	log.Error("Server stopped");

}


func setupLogger(env string) *slog.Logger{
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log;
}

func setupPrettySlog() *slog.Logger{
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}
	handler := opts.NewPrettyHandler(os.Stdout);

	return slog.New(handler);
}