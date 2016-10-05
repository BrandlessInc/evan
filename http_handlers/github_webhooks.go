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
	applications *config.Applications
	githubClient *github.Client
}

func NewGithubEventHandler(applications *config.Applications, githubClient *github.Client) *GithubEventHandler {
	return &GithubEventHandler{
		applications: applications,
		githubClient: githubClient,
	}
}

func (handler *GithubEventHandler) HandleDeploymentEvent(deploymentEvent *github.DeploymentEvent) error {
	app := handler.applications.FindApplicationForGithubRepository(deploymentEvent.Repo)

	target, strategy := app.DeployEnvironment(*deploymentEvent.Deployment.Environment)

	_ = &context.Deployment{
		Application: app,
		Target: target,
		Strategy: strategy,
		Ref: *deploymentEvent.Deployment.Ref,
		GithubClient: handler.githubClient,
	}

	return nil
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

	event := req.Header["X-GitHub-Event"][0]

	if event == "deployment" {
		var deploymentEvent github.DeploymentEvent
		err := json.Unmarshal(body, &deploymentEvent)
		if err != nil {
			respondWithError(res, err)
			return
		}
		respondWithOk(res, "OK")
	} else {
		message := fmt.Sprintf("Cannot handle event: %v", event)
		http.Error(res, message, http.StatusNotImplemented)
	}
}
