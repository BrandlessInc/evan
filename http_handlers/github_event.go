package http_handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Everlane/evan/application"
	"github.com/Everlane/evan/strategy"

	"github.com/google/go-github/github"
)

type GithubEventHandler struct {
	suite        *application.Suite
	githubClient *github.Client
}

func NewGithubEventHandler(suite *application.Suite, githubClient *github.Client) *GithubEventHandler {
	return &GithubEventHandler{
		suite:        suite,
		githubClient: githubClient,
	}
}

func (handler *GithubEventHandler) HandleDeploymentEvent(deploymentEvent *github.DeploymentEvent) error {
	repo := deploymentEvent.Repo
	app := handler.suite.FindApplicationForGithubRepository(repo)

	environment := *deploymentEvent.Deployment.Environment
	target := app.TargetForEnvironment(environment)
	strat := app.StrategyForEnvironment(environment)

	runner := &strategy.Runner{
		Repository:   app.Repository,
		GithubClient: handler.githubClient,

		Strategy: strat,
		Target:   target,

		Ref: *deploymentEvent.Deployment.Ref,
	}

	return runner.Run()
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
