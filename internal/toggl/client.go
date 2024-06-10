package toggl

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
		url:        "https://api.track.toggl.com/api/v9",
		httpClient: &http.Client{},
	}
}

func (c *Client) request(method, path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.url, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.token, "api_token")

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

type TimeEntry struct {
	ID          int    `json:"id"`
	At          string `json:"at"`
	ClientName  string `json:"client_name"`
	Description string `json:"description"`
	ProjectName string `json:"project_name"`
	WorkspaceID int    `json:"workspace_id"`
}

func (c *Client) GetCurrentTimeEntry() (*TimeEntry, error) {
	path := "/me/time_entries/current"
	body, err := c.request(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var data TimeEntry
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) StopTimeEntry(te *TimeEntry) error {
	path := fmt.Sprintf(
		"/workspaces/%d/time_entries/%d/stop",
		te.WorkspaceID,
		te.ID,
	)
	_, err := c.request(http.MethodPatch, path)
	if err != nil {
		return err
	}

	return nil
}
