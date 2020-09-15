package feedrefresh

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func MakeHttpHandler(s Service) http.Handler {
	router := httprouter.New()
	hh := httpHandlers{}

	router.PUT("/articles", hh.refreshArticles)

	return router
}

type httpHandlers struct {
	service Service
}

func (hh httpHandlers) refreshArticles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := hh.service.RefreshArticles(r.Context())
	if err != nil {
		logUnhandledError(errors.Wrap(err, "refresh articles request failed"))
		httpInternalServerError(w)
	}
}

func httpInternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func logUnhandledError(err error) {
	log.Printf("Unhandled error: %s\n", err)
}
