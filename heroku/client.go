package heroku

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	HEADER_ACCEPT       = "application/vnd.heroku+json; version=3"
	HEADER_CONTENT_TYPE = "application/json"
)

type Client struct {
	Token      string
	httpClient *http.Client
}

type SourceBlob struct {
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url,omitempty"`
	Version  string `json:"version,omitempty"`
}

func (c *Client) BaseUrl() string {
	return "https://api.heroku.com"
}

func (c *Client) BuildCreate(appId string, sourceBlob *SourceBlob) error {
	body, err := json.Marshal(map[string]interface{}{
		"source_blob": sourceBlob,
	})
	if err != nil {
		return err
	}

	url := c.BaseUrl() + fmt.Sprintf("/apps/%v/builds", appId)
	_, err = c.MakeRequest("POST", url, &body)

	return nil
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

	return c.httpClient.Do(req)
}
