package redirect

import (
	"context"
	"errors"
	resp "link-service/internal/lib/api/response"
	"link-service/internal/lib/logger/logslog"
	postgres "link-service/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		resURL, err := urlGetter.GetURL(r.Context(), alias)
		if errors.Is(err, postgres.UrlNotFound) {
			log.Info("not found url", "alias: ", alias)
			render.JSON(w, r, resp.Error("not found url"))
			return
		}
		if err != nil {
			log.Info("failed to get url", logslog.Error(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("got url",
			slog.String("url", resURL),
		)
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
