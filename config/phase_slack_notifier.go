package config

import (
    "github.com/nlopes/slack"
)

type SlackNotifierPhase struct {
    Client *slack.Client
    Channel string
    Format func(Deployment) (string, error)
}

func (snp *SlackNotifierPhase) CanPreload() bool {
    return false
}

func (snp *SlackNotifierPhase) CanExecute() (bool, error) {
    return true, nil
}

func (snp *SlackNotifierPhase) IsExecuting() (bool, error) {
    // This executes immediately, so we don't need to poll its status.
    return false, nil
}

func (snp *SlackNotifierPhase) HasExecuted() (bool, error) {
    // Assume it's already happened.
    return true, nil
}

func (snp *SlackNotifierPhase) Execute(deployment Deployment) (ExecuteStatus, error) {
    message, err := snp.Format(deployment)
    if err != nil {
        return ERROR, err
    }

    params := slack.NewPostMessageParameters()
    _, _, err = snp.Client.PostMessage(snp.Channel, message, params)
    if err != nil {
        return ERROR, err
    }

    return DONE, nil
}
