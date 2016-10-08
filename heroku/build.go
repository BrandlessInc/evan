package heroku

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Build struct {
	Id         string      `json:"id"`
	SourceBlob *SourceBlob `json:"source_blob"`
	Status     string      `json:"status"`
}

type SourceBlob struct {
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url"`
	Version  string `json:"version"`
}

type buildCreateRequest struct {
	SourceBlob *SourceBlob `json:"source_blob"`
}

func (c *Client) BuildCreate(appId string, sourceBlob *SourceBlob) (*Build, *http.Response, error) {
	fmt.Printf("sourceBlob: %+v\n", sourceBlob)
	body, err := json.Marshal(&buildCreateRequest{
		SourceBlob: sourceBlob,
	})
	fmt.Printf("BuildCreate body: %v\n", string(body))
	if err != nil {
		return nil, nil, err
	}

	url := c.BaseUrl() + fmt.Sprintf("/apps/%v/builds", appId)
	resp, err := c.MakeRequest("POST", url, &body)
	if err != nil {
		return nil, resp, err
	}

	var build Build
	err = c.readResponseInto(resp, &build)
	if err != nil {
		return nil, resp, err
	}

	return &build, resp, nil
}

func (c *Client) BuildInfo(appId string, id string) (*Build, *http.Response, error) {
	url := c.BaseUrl() + fmt.Sprintf("/apps/%v/builds/%v", appId, id)
	resp, err := c.MakeRequest("GET", url, nil)
	if err != nil {
		return nil, resp, err
	}

	var build Build
	err = c.readResponseInto(resp, &build)
	if err != nil {
		return nil, resp, err
	}

	return &build, resp, nil
}

func (c *Client) readResponse(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 && c.Debug {
		fmt.Printf("[Heroku] Error: code=%v %v\n", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) readResponseInto(resp *http.Response, val interface{}) error {
	body, err := c.readResponse(resp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, val)
	if err != nil {
		return err
	}

	return nil
}
