package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func ReturnJSONData(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(obj)
	return errors.Wrap(err, "can not encode or write data into http response")
}
