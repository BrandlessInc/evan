package heroku

import (
	"bytes"
	"io"
	"net/http"
)

const (
	HEADER_ACCEPT       = "application/vnd.heroku+json; version=3"
	HEADER_CONTENT_TYPE = "application/json"
)

type Client struct {
	Token      string
	Debug      bool
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		httpClient: &http.Client{},
	}
}

func (c *Client) BaseUrl() string {
	return "https://api.heroku.com"
}

func (c *Client) MakeRequest(method, url string, body *[]byte) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(*body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", HEADER_ACCEPT)
	req.Header.Set("Content-Type", HEADER_CONTENT_TYPE)
	req.Header.Set("Authorization", "Bearer " + c.Token)

	return c.httpClient.Do(req)
}
