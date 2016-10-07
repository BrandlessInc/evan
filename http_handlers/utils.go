package http_handlers

import (
	"fmt"
	"net/http"
	"strings"
)

// "Smart" error responder; does a little introspection on the error message
// to try to respond with the most accurate possible HTTP status code.
func respondWithError(res http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if strings.Contains(strings.ToLower(err.Error()), "cannot") {
		status = http.StatusNotImplemented
	}

	http.Error(res, err.Error(), status)
}

func respondWithOk(res http.ResponseWriter, message string) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(200)
	fmt.Fprintln(res, message)
}
