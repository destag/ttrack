package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/destag/ttrack/internal/tasks"
)

type Client struct {
	username   string
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Jira client with basic auth
// You can use username/token or email/token depending on your Jira setup
func NewClient(username, token, baseURL string) *Client {
	return &Client{
		username:   username,
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"` // URL to the issue
	Fields Fields `json:"fields"`
}

func (i Issue) toTask() *tasks.Task {
	return &tasks.Task{
		ID:          i.Key,
		Description: i.Fields.Summary,
	}
}

type Fields struct {
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Created     string    `json:"created"`
	Updated     string    `json:"updated"`
	Status      Status    `json:"status"`
	Priority    Priority  `json:"priority"`
	IssueType   IssueType `json:"issuetype"`
}

type Status struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
}

type IssueType struct {
	Name string `json:"name"`
}

// GetIssue retrieves a Jira issue by its key (e.g., "PROJ-123")
func (c *Client) GetTask(issueKey string) (*tasks.Task, error) {
	path := fmt.Sprintf("/rest/api/2/issue/%s", issueKey)
	body, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	var issue Issue
	err = json.Unmarshal(body, &issue)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue: %w", err)
	}

	return issue.toTask(), nil
}

func (c *Client) ListTasks(query string) ([]*tasks.Task, error) {
	encodedJQL := url.QueryEscape(query)
	path := fmt.Sprintf("/rest/api/2/search?jql=%s", encodedJQL)
	body, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	var resp struct {
		Issues []Issue `json:"issues"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal issues: %w", err)
	}

	var out []*tasks.Task
	for _, issue := range resp.Issues {
		out = append(out, issue.toTask())
	}
	return out, nil
}
