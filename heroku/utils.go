package heroku

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
