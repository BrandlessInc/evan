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

type GithubEventHandler struct {
	Applications *config.Applications
	PreDeployment func(*http.Request, *context.Deployment) error
}

func (handler *GithubEventHandler) HandleDeploymentEvent(req *http.Request, deploymentEvent *github.DeploymentEvent) error {
	app := handler.Applications.FindApplicationForGithubRepository(deploymentEvent.Repo)

	environment := *deploymentEvent.Deployment.Environment
	strategy := app.DeployEnvironment(environment)

	deployment := &context.Deployment{
		Application: app,
		Environment: environment,
		Strategy: strategy,
		Ref: *deploymentEvent.Deployment.Ref,
	}

	if handler.PreDeployment != nil {
		err := handler.PreDeployment(req, deployment)
		if err != nil {
			return err
		}
	}

	return deployment.Run()
}

func respondWithError(res http.ResponseWriter, err error) {
	http.Error(res, fmt.Sprintf("%v", err), http.StatusInternalServerError)
}

func respondWithOk(res http.ResponseWriter, message string) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(200)
	fmt.Fprintln(res, message)
}

func (handler *GithubEventHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(res, err)
		return
	}

	event := req.Header.Get("X-GitHub-Event")

	if event == "deployment" {
		var deploymentEvent github.DeploymentEvent
		err := json.Unmarshal(body, &deploymentEvent)
		if err != nil {
			respondWithError(res, err)
			return
		}
		respondWithOk(res, "OK")
	} else if event == "ping" {
		respondWithOk(res, "PONG")
	} else {
		message := fmt.Sprintf("Cannot handle event: %v", event)
		http.Error(res, message, http.StatusNotImplemented)
	}
}
