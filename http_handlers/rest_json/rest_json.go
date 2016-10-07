package rest_json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/config"
	"github.com/Everlane/evan/context"
)

type CreateDeploymentRequest struct {
	Application string `json:"application"`
	Environment string `json:"environment"`
	Ref         string `json:"ref"`
}

func (cdr *CreateDeploymentRequest) newDeployment(app *config.Application) (*context.Deployment, error) {
	return context.NewDeployment(app.Wrapper(), cdr.Environment, cdr.Ref)
}

type CreateDeploymentHandler struct {
	Applications *config.Applications
	PreDeployment func(*http.Request, common.Deployment) error
}

func (handler *CreateDeploymentHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var deploymentRequest CreateDeploymentRequest
	err := handler.readRequestInto(req, &deploymentRequest)
	if err != nil {
		respondWithError(res, err, http.StatusInternalServerError)
		return
	}

	app := handler.Applications.FindApplicationByName(deploymentRequest.Application)
	if app == nil {
		err = fmt.Errorf("Application not found: '%v'", deploymentRequest.Application)
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

	if handler.PreDeployment != nil {
		err = handler.PreDeployment(req, deployment)
		if err != nil {
			respondWithError(res, err, http.StatusInternalServerError)
			return
		}
	}

	err = deployment.CheckPreconditions()
	if err != nil {
		respondWithError(res, err, http.StatusPreconditionFailed)
		return
	}

	humanDescription := fmt.Sprintf(
		"%v to %v for %v",
		deploymentRequest.Ref,
		deploymentRequest.Environment,
		deploymentRequest.Application,
	)

	// Start the party!
	go func() {
		err := deployment.RunPhases()
		if err != nil {
			fmt.Printf("Error deploying %v: %v", humanDescription, err)
		}
	}()

	message := fmt.Sprintf("Deploying %v", humanDescription)
	respondWithOk(res, message)
}

func (handler *CreateDeploymentHandler) readRequestInto(req *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, val)
}
