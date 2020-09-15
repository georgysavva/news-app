package feed

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func MakeHttpHandler(s Service) http.Handler {
	router := httprouter.New()
	hh := httpHandlers{}

	router.GET("/articles", hh.getArticles)
	router.GET("/articles/:id", hh.getArticle)
	router.GET("/categories", hh.getCategories)
	router.GET("/providers", hh.getProviders)

	return router
}

type httpHandlers struct {
	service Service
}

func (hh httpHandlers) getArticles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryParams := r.URL.Query()
	categories := queryParams["categories"]
	providers := queryParams["providers"]
	articles, err := hh.service.GetArticles(r.Context(), WithCategories(categories), WithProvides(providers))
	if err != nil {
		logUnhandledError(errors.Wrap(err, "get articles request failed"))
		httpInternalServerError(w)
	}
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		log.Println("Can not encode or write articles data into http response")
	}
}

func (hh httpHandlers) getArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	articleID := params.ByName("id")

	a, err := hh.service.GetArticle(r.Context(), articleID)
	if err != nil {
		logUnhandledError(errors.Wrap(err, "get article request failed"))
		httpInternalServerError(w)
	}
	if a == nil {
		http.Error(w, fmt.Sprintf("Article %s not found", articleID), http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(a); err != nil {
		log.Println("Can not encode or write article data into http response")
	}
}

func (hh httpHandlers) getCategories(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	categories, err := hh.service.GetCategories(r.Context())
	if err != nil {
		logUnhandledError(errors.Wrap(err, "get categories request failed"))
		httpInternalServerError(w)
	}
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Println("Can not encode or write categories list into http response")
	}
}

func (hh httpHandlers) getProviders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	providers, err := hh.service.GetProviders(r.Context())
	if err != nil {
		logUnhandledError(errors.Wrap(err, "get categories failed"))
		httpInternalServerError(w)
	}
	if err := json.NewEncoder(w).Encode(providers); err != nil {
		log.Println("Can not encode or write providers list into http response")
	}
}

func httpInternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func logUnhandledError(err error) {
	log.Printf("Unhandled error: %s\n", err)
}
