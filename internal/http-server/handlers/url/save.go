package save

import (
	"errors"
	"github.com/20grizz03/restApiURLShortener/internal/db"
	"github.com/20grizz03/restApiURLShortener/internal/lib/api/response"
	"github.com/20grizz03/restApiURLShortener/internal/lib/logger/sl"
	"github.com/20grizz03/restApiURLShortener/internal/lib/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

const aliasLegth = 6

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURl(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, URLSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, response.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidatorError(validateErr))

			return

		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLegth)
		}

		id, err := URLSaver.SaveURl(req.URL, alias)
		if errors.Is(err, db.ErrUrlExists) {
			log.Info("url already exists", slog.String("alias", alias))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to save url"))

			return
		}

		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
