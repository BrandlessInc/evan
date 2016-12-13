package phases

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/heroku"
)

const TOKEN_FLAG string = "heroku.token"

// Creates a build via the Heroku API with the deployment's ref/SHA1. The
// build will use a generated tarball of the repository from GitHub (see its
// Preload method for details).
type HerokuBuildPhase struct {
	// Will use this client if no client is passed through build flags.
	DefaultClient *heroku.Client
	AppId         string
}

type HerokuBuildPhaseContext struct {
	herokuClient *heroku.Client
	sha1         string
	tarballUrl   string
}

func (hbp *HerokuBuildPhase) CanPreload() bool {
	return true
}

func (hbp *HerokuBuildPhase) Preload(deployment common.Deployment) (interface{}, error) {
	herokuClient := hbp.DefaultClient
	if deployment.HasFlag(TOKEN_FLAG) {
		token := deployment.Flag(TOKEN_FLAG).(string)
		herokuClient = heroku.NewClient(token)
	}
	if herokuClient == nil {
		return nil, fmt.Errorf("No Heroku Platform API client found")
	}

	githubRepo, err := common.NewGithubRepositoryFromDeployment(deployment)
	if err != nil {
		return nil, err
	}

	sha1, err := githubRepo.GetCommitSHA1(deployment.Ref())
	if err != nil {
		return nil, err
	}

	tarballUrl, err := githubRepo.GetArchiveLink(common.Tarball, deployment.Ref())
	if err != nil {
		return nil, err
	}

	return &HerokuBuildPhaseContext{
		herokuClient: herokuClient,
		sha1:         sha1,
		tarballUrl:   tarballUrl,
	}, nil
}

func (hbp *HerokuBuildPhase) Execute(deployment common.Deployment, data interface{}) error {
	context := data.(*HerokuBuildPhaseContext)

	build, err := hbp.createBuild(deployment, context)
	if err != nil {
		return err
	}

	build, err = hbp.PollBuildStatus(build, context)
	if err != nil {
		return err
	}
	if build.Status == "failed" {
		return fmt.Errorf("Build %v failed", build.Id)
	}

	// Heroku will have automatically deployed that new build as a release, so
	// we can consider this phase done.
	return nil
}

func (hbp *HerokuBuildPhase) createBuild(deployment common.Deployment, context *HerokuBuildPhaseContext) (*heroku.Build, error) {
	build, resp, err := context.herokuClient.BuildCreate(hbp.AppId, &heroku.SourceBlob{
		Url:     context.tarballUrl,
		Version: context.sha1,
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Received error from Heroku when creating build: %v\n", resp)
		return nil, fmt.Errorf("Non-creation status code (expected %v, got %v)", http.StatusCreated, resp.StatusCode)
	}

	return build, nil
}

func (hbp *HerokuBuildPhase) PollBuildStatus(sourceBuild *heroku.Build, context *HerokuBuildPhaseContext) (*heroku.Build, error) {
	for true {
		build, _, err := context.herokuClient.BuildInfo(hbp.AppId, sourceBuild.Id)
		if err != nil {
			return nil, err
		}

		if build.Status == "pending" {
			time.Sleep(10 * time.Second)
			continue
		}

		return build, nil
	}

	panic("Unreachable!")
}
