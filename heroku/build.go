package heroku

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SourceBlob struct {
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url"`
	Version  string `json:"url"`
}

func (c *Client) BuildCreate(appId string, sourceBlob *SourceBlob) (*http.Response, error) {
	body, err := json.Marshal(map[string]interface{}{
		"source_blob": sourceBlob,
	})
	if err != nil {
		return nil, err
	}

	url := c.BaseUrl() + fmt.Sprintf("/apps/%v/builds", appId)
	return c.MakeRequest("POST", url, &body)
}
