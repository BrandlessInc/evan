package slack

// Adapted from:
//   https://github.com/ashwanthkumar/slack-go-webhook
//   https://github.com/huguesalary/slack-go/blob/master/slack.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Attachment struct {
	Color string `json:"color,omitempty"`
	Text  string `json:"text,omitempty"`
}

type Payload struct {
	Username    string       `json:"username,omitempty"`
	IconUrl     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type WebhookError struct {
	StatusCode int
	Body       string
}

func (err *WebhookError) Error() string {
	body := strings.TrimSpace(err.Body)
	return fmt.Sprintf("WebhookError(%d): %s", err.StatusCode, body)
}

func Send(webhookUrl string, payload Payload) error {
	payloadBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(payloadBody)

	resp, err := http.Post(webhookUrl, "application/json", buf)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		} else {
			return &WebhookError{
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}
		}
	}

	return nil
}
