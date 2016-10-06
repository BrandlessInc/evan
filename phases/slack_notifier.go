package phases

import (
	"github.com/Everlane/evan/common"

	"github.com/nlopes/slack"
)

type SlackNotifierPhase struct {
	Client  *slack.Client
	Channel string
	Format  func(common.Deployment) (string, error)
}

func (snp *SlackNotifierPhase) CanPreload() bool {
	return false
}

func (snp *SlackNotifierPhase) HasExecuted(deployment common.Deployment) (bool, error) {
	return false, nil
}

func (snp *SlackNotifierPhase) Execute(deployment common.Deployment) (common.ExecuteStatus, error) {
	message, err := snp.Format(deployment)
	if err != nil {
		return common.PHASE_ERROR, err
	}

	// If the `Format` function returned an empty strings that means we
	// shouldn't send a message to Slack.
	if message == "" {
		return common.PHASE_DONE, nil
	}

	params := slack.NewPostMessageParameters()
	_, _, err = snp.Client.PostMessage(snp.Channel, message, params)
	if err != nil {
		return common.PHASE_ERROR, err
	}

	return common.PHASE_DONE, nil
}
