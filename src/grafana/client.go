package grafana

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL    string
	APIKey     string
	httpClient *http.Client
}

func NewClient(baseURL string, apiKey string) (c *Client) {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		httpClient: new(http.Client),
	}
}

func (c *Client) request(method string, endpoint string, body []byte) ([]byte, error) {
	url := c.BaseURL + "/api/" + endpoint

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	authHeader := fmt.Sprintf("Bearer %s", c.APIKey)
	req.Header.Add("Authorization", authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			err = fmt.Errorf("%s not found (404)", url)
		} else {
			err = fmt.Errorf(
				"Unknown error: %d; body: %s",
				resp.StatusCode,
				string(respBody),
			)
		}
	}

	return respBody, err
}