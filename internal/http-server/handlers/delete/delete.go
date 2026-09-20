package delete

import (
	"context"
	resp "link-service/internal/lib/api/response"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type DelURL interface {
	DeleteURL(ctx context.Context, id int64) error
}

func New(log *slog.Logger, delURL DelURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("cannot convertation id")
			render.JSON(w, r, resp.Error("inernal server"))
			return
		}
		if id <= 0 {
			log.Error("invalid value id")
			render.JSON(w, r, resp.Error("inernal server"))
			return
		}
		err = delURL.DeleteURL(r.Context(), id)
		if err != nil {
			log.Error("id is not delete")
			render.JSON(w, r, resp.Error("inernal server"))
			return
		}
		log.Info("url successfully deleted", slog.Int64("id", id))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.OK());

	}
}
