package http_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(res http.ResponseWriter, err error, code int) {
	body, err := json.Marshal(map[string]interface{}{
		"message": err.Error(),
	})
	if err != nil {
		panic(err) // Something real bad happened
	}

	http.Error(res, string(body), code)
}

func respondWithOk(res http.ResponseWriter, message string) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(200)
	fmt.Fprintln(res, message)
}
