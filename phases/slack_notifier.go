package phases

import (
	"github.com/Everlane/evan/common"

	"github.com/nlopes/slack"
)

type SlackNotifierPhase struct {
	Client  *slack.Client
	Channel string
	Format  func(common.Deployment) (*string, *slack.PostMessageParameters, error)
}

func (snp *SlackNotifierPhase) CanPreload() bool {
	return false
}

func (snp *SlackNotifierPhase) Execute(deployment common.Deployment, _ interface{}) error {
	message, params, err := snp.Format(deployment)
	if err != nil {
		return err
	}

	// Don't send a message to Slack if the format function didn't return
	// a message to send
	if message == nil {
		return nil
	}

	if params == nil {
		defaultParams := slack.NewPostMessageParameters()
		params = &defaultParams
	}

	_, _, err = snp.Client.PostMessage(snp.Channel, *message, *params)
	if err != nil {
		return err
	}

	return nil
}
