package phases

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/heroku"
)

type HerokuBuildPhase struct {
	Client *heroku.Client
	AppId  string
}

type HerokuBuildPhaseContext struct {
	TarballUrl string
}

func (hbp *HerokuBuildPhase) CanPreload() bool {
	return true
}

func (hbp *HerokuBuildPhase) Preload(deployment common.Deployment) (interface{}, error) {
	githubRepo := common.NewGithubRepositoryFromDeployment(deployment)
	tarballUrl, err := githubRepo.GetArchiveLink(common.Tarball)
	if err != nil {
		return nil, err
	}

	return &HerokuBuildPhaseContext{
		TarballUrl: tarballUrl,
	}, nil
}

func (hbp *HerokuBuildPhase) Execute(deployment common.Deployment, data interface{}) error {
	context := data.(*HerokuBuildPhaseContext)

	build, err := hbp.createBuild(deployment, context)
	if err != nil {
		return err
	}

	build, err = hbp.PollBuildStatus(build)
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
	build, resp, err := hbp.Client.BuildCreate(hbp.AppId, &heroku.SourceBlob{
		Url:     context.TarballUrl,
		Version: deployment.Ref(),
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

func (hbp *HerokuBuildPhase) PollBuildStatus(sourceBuild *heroku.Build) (*heroku.Build, error) {
	for true {
		build, _, err := hbp.Client.BuildInfo(hbp.AppId, sourceBuild.Id)
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
