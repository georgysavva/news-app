package feed

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/georgysavva/news-app/pkg/httputil"
)

func MakeHTTPHandler(s Service) http.Handler {
	router := mux.NewRouter()
	hh := httpHandlers{service: s}

	router.HandleFunc("/articles", hh.getArticles).Methods("GET")
	router.HandleFunc("/articles/{id}", hh.getArticle).Methods("GET")
	router.HandleFunc("/categories", hh.getCategories).Methods("GET")
	router.HandleFunc("/providers", hh.getProviders).Methods("GET")

	return router
}

type httpHandlers struct {
	service Service
}

func (hh httpHandlers) getArticles(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	categories := queryParams["categories"]
	providers := queryParams["providers"]
	articles, err := hh.service.GetArticles(r.Context(), WithCategories(categories), WithProvides(providers))
	if err != nil {
		logError(errors.Wrap(err, "get articles request failed"))
		httputil.InternalServerError(w)
		return
	}
	if err := httputil.ReturnJSONData(w, articles); err != nil {
		logError(errors.Wrap(err, "can not return articles list"))
	}
}

func (hh httpHandlers) getArticle(w http.ResponseWriter, r *http.Request) {
	articleID := mux.Vars(r)["id"]

	a, err := hh.service.GetArticle(r.Context(), articleID)
	if err != nil {
		logError(errors.Wrap(err, "get article request failed"))
		httputil.InternalServerError(w)
		return
	}
	if a == nil {
		http.Error(w, fmt.Sprintf("Article %s not found", articleID), http.StatusNotFound)
		return
	}
	if err := httputil.ReturnJSONData(w, a); err != nil {
		logError(errors.Wrap(err, "can not return article object"))
	}
}

func (hh httpHandlers) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := hh.service.GetCategories(r.Context())
	if err != nil {
		logError(errors.Wrap(err, "get categories request failed"))
		httputil.InternalServerError(w)
		return
	}
	if err := httputil.ReturnJSONData(w, categories); err != nil {
		logError(errors.Wrap(err, "can not return categories list"))
	}
}

func (hh httpHandlers) getProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := hh.service.GetProviders(r.Context())
	if err != nil {
		logError(errors.Wrap(err, "get providers request failed"))
		httputil.InternalServerError(w)
		return
	}
	if err := httputil.ReturnJSONData(w, providers); err != nil {
		logError(errors.Wrap(err, "can not return providers list"))
	}
}

func logError(err error) {
	log.Printf("Unhandled error: %s\n", err)
}
