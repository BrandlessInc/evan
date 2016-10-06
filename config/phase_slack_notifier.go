package config

import (
	"github.com/nlopes/slack"
)

type SlackNotifierPhase struct {
	Client  *slack.Client
	Channel string
	Format  func(Deployment) (string, error)
}

func (snp *SlackNotifierPhase) CanPreload() bool {
	return false
}

func (snp *SlackNotifierPhase) HasExecuted(deployment Deployment) (bool, error) {
	return false, nil
}

func (snp *SlackNotifierPhase) Execute(deployment Deployment) (ExecuteStatus, error) {
	message, err := snp.Format(deployment)
	if err != nil {
		return ERROR, err
	}

	// If the `Format` function returned an empty strings that means we
	// shouldn't send a message to Slack.
	if message == "" {
		return DONE, nil
	}

	params := slack.NewPostMessageParameters()
	_, _, err = snp.Client.PostMessage(snp.Channel, message, params)
	if err != nil {
		return ERROR, err
	}

	return DONE, nil
}
