package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// readBody into json struct
func readBody(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return fmt.Errorf("body read: %w", err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}

type jsonError struct {
	Error string `json:"error"`
}

// writeError formatted in json
func writeError(w http.ResponseWriter, err error, statusCode int) {
	writeResponse(w, &jsonError{Error: err.Error()}, statusCode)
}

// writeResponse formatted in json
func writeResponse(w http.ResponseWriter, v interface{}, statusCode int) {
	resBody, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(resBody)
}
