package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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

func (c *Client) request(method, path string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.url, path)
	req, err := http.NewRequest(method, url, body)
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
		s, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println("error body:", string(s), time.Now().UTC().Format(time.TimeOnly))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

type TimeEntry struct {
	ID          int       `json:"id"`
	At          string    `json:"at"`
	Start       time.Time `json:"start"`
	ClientName  string    `json:"client_name"`
	Description string    `json:"description"`
	ProjectName string    `json:"project_name"`
	WorkspaceID int       `json:"workspace_id"`
}

func (c *Client) GetCurrentTimeEntry() (*TimeEntry, error) {
	path := "/me/time_entries/current"
	body, err := c.request(http.MethodGet, path, nil)
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
	_, err := c.request(http.MethodPatch, path, nil)
	if err != nil {
		return err
	}

	return nil
}
