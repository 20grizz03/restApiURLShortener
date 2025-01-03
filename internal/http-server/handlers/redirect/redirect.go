package redirect

import (
	"errors"
	"github.com/20grizz03/restApiURLShortener/internal/db"
	"github.com/20grizz03/restApiURLShortener/internal/lib/api/response"
	"github.com/20grizz03/restApiURLShortener/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.50.3 --name=URLGetter
type URLGetter interface {
	GetUrl(alias string) (string, error)
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

			render.JSON(w, r, response.Error("Not found"))

			return
		}

		resURL, err := urlGetter.GetUrl(alias)
		if errors.Is(err, db.ErrUrlNotFound) {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("Not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal server error"))

			return

		}

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}