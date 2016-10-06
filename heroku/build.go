package heroku

import (
	"bytes"
	"io"
	"net/http"
)

type SourceBlob struct {
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url"`
	Version  string `json:"url"`
}

type Client struct {
	Token string
	httpClient *http.Client
}

func (c *Client) BaseUrl() string {
	return "https://api.heroku.com"
}

func (c *Client) MakeRequest(method, url string, body *[]byte) error {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(*body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return err
	}

	_, err = c.httpClient.Do(req)
	return err
}
