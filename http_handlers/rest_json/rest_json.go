package rest_json

import (
	"encoding/json"
	"fmt"
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

	fmt.Printf("deploymentRequest: %+v", deploymentRequest)

	app := handler.Applications.FindApplicationByName(deploymentRequest.application)
	if app == nil {
		err = fmt.Errorf("Application not found: '%v'", deploymentRequest.application)
		respondWithError(res, err, http.StatusNotFound)
		return
	}

	deployment, err := deploymentRequest.newDeployment(app)
	if err != nil {
		respondWithError(res, err, http.StatusNotFound)
		return
	}

	err = deployment.SetStoreAndSave(handler.Applications.Store)
	if err != nil {
		respondWithError(res, err, http.StatusInternalServerError)
		return
	}

	err = deployment.CheckPreconditions()
	if err != nil {
		respondWithError(res, err, http.StatusPreconditionFailed)
		return
	}

	// Start the party!
	go deployment.RunPhases()

	message := fmt.Sprintf(
		"Deploying %v to %v for %v",
		deploymentRequest.ref,
		deploymentRequest.environment,
		deploymentRequest.application,
	)
	respondWithOk(res, message)
}

func (handler *CreateDeploymentHandler) readRequestInto(req *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &val)
}
