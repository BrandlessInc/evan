package phases

import (
    "fmt"
    "net/http"

    "github.com/Everlane/evan/common"
    "github.com/Everlane/evan/heroku"
)

type HerokuBuildPhase struct {
    Client *heroku.Client
    AppId string
}

func (hbp *HerokuBuildPhase) CanPreload() bool {
    return false
}

func (hbp *HerokuBuildPhase) Execute(deployment common.Deployment) error {
    err := hbp.createBuild(deployment)
    if err != nil {
        return err
    }

    return nil
}

func (hbp *HerokuBuildPhase) createBuild(deployment common.Deployment) error {
    githubRepo := common.NewGithubRepositoryFromDeployment(deployment)
    tarballUrl, err := githubRepo.GetArchiveLink(common.Tarball)
    if err != nil {
        return err
    }

    resp, err := hbp.Client.BuildCreate(hbp.AppId, &heroku.SourceBlob{
        Url: tarballUrl,
        Version: deployment.Ref(),
    })
    if err != nil {
        return err
    }
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("Non-creation status code (expected %v, got %v)", http.StatusCreated, resp.StatusCode)
    }

    return nil
}
