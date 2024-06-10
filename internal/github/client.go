package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	token      string
	url        string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		url:        "https://api.github.com",
		httpClient: &http.Client{},
	}
}

func (c *Client) get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.url, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

type Issue struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *Client) GetIssue(project string, id int) (*Issue, error) {
	path := fmt.Sprintf("/repos/%s/issues/%d", project, id)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data Issue
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
