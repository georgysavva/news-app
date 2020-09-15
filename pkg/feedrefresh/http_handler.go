package feedrefresh

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func MakeHttpHandler(s Service) http.Handler {
	router := mux.NewRouter()
	hh := httpHandlers{service: s}

	router.HandleFunc("/", hh.refreshFeed)

	return router
}

type httpHandlers struct {
	service Service
}

func (hh httpHandlers) refreshFeed(w http.ResponseWriter, r *http.Request) {
	err := hh.service.RefreshFeed(r.Context())
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
