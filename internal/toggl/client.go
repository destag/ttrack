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

func (c *Client) doRequest(method, path string, body io.Reader, out any) error {
	url := fmt.Sprintf("%s%s", c.url, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.token, "api_token")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		s, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println("error body:", string(s), time.Now().UTC().Format(time.TimeOnly))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if out == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(&out)
}

type CurrentUser struct {
	ID                 int `json:"id"`
	DefaultWorkspaceID int `json:"default_workspace_id"`
}

func (c *Client) GetUserInfo() (*CurrentUser, error) {
	var data CurrentUser
	path := "/me"
	err := c.doRequest(http.MethodGet, path, nil, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
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
	var data TimeEntry
	path := "/me/time_entries/current"
	err := c.doRequest(http.MethodGet, path, nil, &data)
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
	return c.doRequest(http.MethodPatch, path, nil, nil)
}

func (c *Client) GetTimeEntries() ([]*TimeEntry, error) {
	var data []*TimeEntry
	path := "/me/time_entries?meta=true"
	err := c.doRequest(http.MethodGet, path, nil, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type StartTimeEntryParams struct {
	Project     int    `json:"project_id"`
	Description string `json:"description"`
	Start       string `json:"start"`
	CreatedWith string `json:"created_with"`
	WorkspaceID int    `json:"workspace_id"`
	Duration    int    `json:"duration"`
}

func (c *Client) StartTimeEntry(workspaceID int, title string, project string) error {
	te, err := c.GetCurrentTimeEntry()
	if err != nil {
		return err
	}

	if te.WorkspaceID != 0 {
		return errors.New("time tracker already started")
	}

	p, err := c.GetProject(workspaceID, project)
	if err != nil {
		return err
	}

	data := StartTimeEntryParams{
		Project:     p.ID,
		Description: title,
		Start:       time.Now().UTC().Format(time.RFC3339),
		CreatedWith: "ttrack",
		WorkspaceID: workspaceID,
		Duration:    -1,
	}

	bs, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(bs)

	path := fmt.Sprintf("/workspaces/%d/time_entries", workspaceID)
	return c.doRequest(http.MethodPost, path, reader, nil)
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetProject(workspaceID int, name string) (*Project, error) {
	var data []Project
	path := fmt.Sprintf("/workspaces/%d/projects", workspaceID)
	err := c.doRequest(http.MethodGet, path, nil, &data)
	if err != nil {
		return nil, err
	}

	for _, p := range data {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, errors.New("project not found")
}
