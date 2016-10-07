package http_handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Everlane/evan/config"
	"github.com/Everlane/evan/context"
)

type CreateDeploymentRequest struct {
	application string `json:"application"`
	environment string `json:"environment"`
	ref         string `json:"ref"`
}

func (cdr *CreateDeploymentRequest) newDeployment(app *config.Application) (*context.Deployment, error) {
	return context.NewDeployment(app.Wrapper(), cdr.environment, cdr.ref)
}

type CreateDeploymentHandler struct {
	Applications *config.Applications
}

func (handler *CreateDeploymentHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var deploymentRequest CreateDeploymentRequest
	err := handler.readRequestInto(req, &deploymentRequest)
	if err != nil {
		respondWithError(res, err, http.StatusInternalServerError)
		return
	}
}

func (handler *CreateDeploymentHandler) readRequestInto(req *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, val)
}
