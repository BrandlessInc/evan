package slack

// Adapted from:
//   https://github.com/ashwanthkumar/slack-go-webhook

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Fallback   *string  `json:"fallback"`
	Color      *string  `json:"color"`
	PreText    *string  `json:"pretext"`
	AuthorName *string  `json:"author_name"`
	AuthorLink *string  `json:"author_link"`
	AuthorIcon *string  `json:"author_icon"`
	Title      *string  `json:"title"`
	TitleLink  *string  `json:"title_link"`
	Text       *string  `json:"text"`
	ImageUrl   *string  `json:"image_url"`
	Fields     []*Field `json:"fields"`
	Footer     *string  `json:"footer"`
	FooterIcon *string  `json:"footer_icon"`
}

type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconUrl     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("Incorrect token (redirection)")
}

func Send(webhookUrl string, proxy string, payload Payload) error {
	request := gorequest.New().Proxy(proxy)
	res, _, errs := request.
		Post(webhookUrl).
		RedirectPolicy(redirectPolicyFunc).
		Send(payload).
		End()

	if errs != nil {
		return errs[0]
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("Error sending message (status: %v)", res.Status)
	}

	return nil
}
