package http_handlers

import (
	"encoding/json"
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

func (handler *GithubEventHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	event := req.Header["X-GitHub-Event"][0]

	if event == "deployment" {
		var deploymentEvent github.DeploymentEvent
		err := json.Unmarshal(body, &deploymentEvent)
		if err != nil {
			panic(err)
		}
	}
}
