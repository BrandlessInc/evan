package rest_json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/config"
	"github.com/Everlane/evan/context"
)

type CreateDeploymentRequest struct {
	Application string                 `json:"application"`
	Environment string                 `json:"environment"`
	Ref         string                 `json:"ref"`
	Flags       map[string]interface{} `json:"flags"`
}

func (cdr *CreateDeploymentRequest) newDeployment(app *config.Application) (*context.Deployment, error) {
	return context.NewDeployment(app.Wrapper(), cdr.Environment, cdr.Ref, cdr.Flags)
}

type CreateDeploymentHandler struct {
	Applications  *config.Applications
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

	hasActiveDeployment, err := handler.Applications.Store.HasActiveDeployment(app.Wrapper(), deploymentRequest.Environment)
	if err != nil {
		respondWithError(res, err, http.StatusInternalServerError)
		return
	}
	if hasActiveDeployment {
		respondWithError(res, fmt.Errorf("Deployment in progress"), http.StatusLocked)
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

	humanDescription := fmt.Sprintf(
		"%v to %v for %v",
		deploymentRequest.Ref,
		deploymentRequest.Environment,
		deploymentRequest.Application,
	)

	// Start the party!
	go func() {
		err := deployment.Execute()
		if err != nil {
			fmt.Printf("Error deploying %v: %v\n", humanDescription, err)
		}
	}()

	body, err := json.Marshal(map[string]interface{}{
		"message":    fmt.Sprintf("Deploying %v", humanDescription),
		"deployment": DeploymentAsJSON(deployment),
	})
	if err != nil {
		panic(err)
	}
	respondWith(res, body, http.StatusCreated)
}

func DeploymentStatusAsJSON(status common.DeploymentStatus) map[string]interface{} {
	json := map[string]interface{}{
		"state": status.State.String(),
		"phase": nil,
		"error": nil,
	}
	if status.Phase != nil {
		json["phase"] = reflect.Indirect(reflect.ValueOf(status.Phase)).Type().Name()
	}
	if status.Error != nil {
		json["error"] = status.Error.Error()
	}
	return json
}

func DeploymentAsJSON(deployment common.Deployment) map[string]interface{} {
	return map[string]interface{}{
		"uuid":        deployment.UUID(),
		"application": deployment.Application().Name(),
		"environment": deployment.Environment(),
		"ref":         deployment.Ref(),
		"sha1":        deployment.SHA1(),
		"status":      DeploymentStatusAsJSON(deployment.Status()),
	}
}

func (handler *CreateDeploymentHandler) readRequestInto(req *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, val)
}

type DeploymentsStatusHandler struct {
	Applications *config.Applications
}

func (handler *DeploymentsStatusHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	store := handler.Applications.Store
	applicationsStatuses := make(map[string]map[string]interface{})

	for _, app := range handler.Applications.List {
		environmentsStatuses := make(map[string]interface{})

		for _, env := range app.Environments {
			deployment, _ := store.FindDeployment(app.Wrapper(), env)

			if deployment != nil {
				environmentsStatuses[env] = DeploymentAsJSON(deployment)
			} else {
				environmentsStatuses[env] = nil
			}
		}

		applicationsStatuses[app.Name] = environmentsStatuses
	}

	body, err := json.Marshal(map[string]interface{}{
		"applications": applicationsStatuses,
	})
	if err != nil {
		respondWithError(res, err, http.StatusInternalServerError)
		return
	}
	respondWith(res, body, http.StatusOK)
}
