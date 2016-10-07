package http_handlers

import (
	// "io/ioutil"
	"net/http"

	"github.com/Everlane/evan/config"
)

type CreateDeploymentRequest struct {
	Ref string `json:"ref"`
}

type CreateDeploymentHandler struct {
	Applications *config.Applications
}

func (handler *CreateDeploymentHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// body, err := ioutil.ReadAll(req.Body)
}
