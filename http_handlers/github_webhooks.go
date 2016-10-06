package http_handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Everlane/evan/config"
	"github.com/Everlane/evan/context"

	"github.com/google/go-github/github"
)

func createDeployment(app *config.Application, environment string, ref string) *context.Deployment {
	strategy := app.DeployEnvironment(environment)

	return context.NewDeployment(app, environment, strategy, ref)
}

type GithubEventHandler struct {
	Applications  *config.Applications
	PreDeployment func(*http.Request, *context.Deployment) error
	PreDeploymentStatus func(*http.Request, *context.Deployment) error
}

func (handler *GithubEventHandler) HandleDeploymentEvent(req *http.Request, body []byte) error {
	var deploymentEvent github.DeploymentEvent
	err := json.Unmarshal(body, &deploymentEvent)
	if err != nil {
		return err
	}

	app := handler.Applications.FindApplicationForGithubRepository(deploymentEvent.Repo)
	environment := *deploymentEvent.Deployment.Environment
	ref := *deploymentEvent.Deployment.Ref

	deployment := createDeployment(app, environment, ref)

	if handler.PreDeployment != nil {
		err := handler.PreDeployment(req, deployment)
		if err != nil {
			return err
		}
	}

	return nil // deployment.Run()
}

func (handler *GithubEventHandler) HandleDeploymentStatusEvent(req *http.Request, body []byte) error {
	var deploymentStatusEvent github.DeploymentStatusEvent
	err := json.Unmarshal(body, &deploymentStatusEvent)
	if err != nil {
		return err
	}

	app := handler.Applications.FindApplicationForGithubRepository(deploymentStatusEvent.Repo)
	environment := *deploymentStatusEvent.Deployment.Environment
	ref := *deploymentStatusEvent.Deployment.Ref

	deployment := createDeployment(app, environment, ref)

	if handler.PreDeploymentStatus != nil {
		err := handler.PreDeploymentStatus(req, deployment)
		if err != nil {
			return err
		}
	}

	return nil // deployment.Run()
}

func (handler *GithubEventHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(res, err)
		return
	}

	event := req.Header.Get("X-GitHub-Event")

	err = handler.HandleEvent(req, event, body)
	if err != nil {
		respondWithError(res, err)
	} else {
		respondWithOk(res, "OK")
	}
}

func (handler *GithubEventHandler) HandleEvent(req *http.Request, event string, body []byte) error {
	switch event {
	case "deployment":
		return handler.HandleDeploymentEvent(req, body)
	case "deployment_status":
		return handler.HandleDeploymentStatusEvent(req, body)
	case "ping":
		return nil
	default:
		return fmt.Errorf("Cannot handle event: %v", event)
	}
}
