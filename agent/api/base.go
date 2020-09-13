package api

import (
	"encoding/json"
	"net/http"
)

func handleInternalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func handleUnauthorized(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized),
		http.StatusUnauthorized)
}

// response status code will be written automatically if there is an error
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		handleInternalServerError(w)
		return err
	}

	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}
