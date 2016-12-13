package rest_json

import (
	"encoding/json"
	"net/http"
)

func bodyWithMessage(message string) []byte {
	body, err := json.Marshal(map[string]interface{}{
		"message": message,
	})
	if err != nil {
		panic(err) // Something real bad happened
	}
	return body
}

func respondWithError(res http.ResponseWriter, err error, code int) {
	body := bodyWithMessage(err.Error())
	respondWith(res, body, code)
}

func respondWithOk(res http.ResponseWriter, message string) {
	body := bodyWithMessage(message)
	respondWith(res, body, http.StatusOK)
}

func respondWith(res http.ResponseWriter, body []byte, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(body)
}
