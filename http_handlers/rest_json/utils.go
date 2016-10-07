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
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	body := bodyWithMessage(err.Error())
	res.Write(body)
}

func respondWithOk(res http.ResponseWriter, message string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(bodyWithMessage(message))
}
