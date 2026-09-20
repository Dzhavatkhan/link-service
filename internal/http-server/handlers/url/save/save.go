package save

import (
	"context"
	"errors"
	"link-service/internal/http-server/handlers/structs"
	resp "link-service/internal/lib/api/response"
	"link-service/internal/lib/logger/logslog"
	"link-service/internal/lib/random"
	postgres "link-service/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v3
type URLSaver interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req structs.Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", logslog.Error(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)

			log.Error("invalid request", logslog.Error(err))

			render.JSON(w, r, resp.ValidationError(validateError))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias, err = random.NewRandomString(aliasLength)
			if err != nil {
				log.Error("invalid random", logslog.Error(err))
				return 
			}
		}

		id, err := urlSaver.SaveURL(r.Context(), req.URL, alias)
		if errors.Is(err, postgres.ErrUrlExists){
			log.Error("url already exists", logslog.Error(err))
			render.JSON(w, r, resp.Error("url already exists"))		
			return 	
		}
		if err != nil {
			log.Error("FAILED TO ADD URL", logslog.Error(err))
			render.JSON(w, r, resp.Error("FAILED TO ADD URL"))
			return 			
		}
		log.Info("url added", slog.Int64("id", id))
		render.JSON(w,r,structs.Response{
			Response: resp.OK(),
			Alias: alias,
		});
	}
}
