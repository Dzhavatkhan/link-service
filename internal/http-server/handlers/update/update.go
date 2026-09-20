package update

import (
	"context"
	"link-service/internal/http-server/handlers/structs"
	resp "link-service/internal/lib/api/response"
	"link-service/internal/lib/logger/logslog"
	"link-service/internal/lib/random"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 6;

type UpdUrl interface {
	UpdateURL(ctx context.Context, id int64, urlToUpd string, aliasToUpd string) error
}

func New(log *slog.Logger, updurl UpdUrl) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.update.New"
		
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

		idStr := chi.URLParam(r, "id");
		if idStr == ""{
			log.Error("id is empty")
			render.JSON(w, r, resp.Error("Internal server error"))
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64);
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
		err = updurl.UpdateURL(r.Context(), id, req.URL, alias);
		if err != nil {
			log.Error("FAILED TO UPDATE URL", logslog.Error(err))
			render.JSON(w, r, resp.Error("FAILED TO UPDATE URL"))
			return 			
		}
		log.Info("url added", slog.Int64("id", id))
		render.JSON(w,r,structs.Response{
			Response: resp.OK(),
			Alias: alias,
		});
	}
}