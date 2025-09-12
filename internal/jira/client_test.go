package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func testServer(t *testing.T, response string) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))

	return server
}

func TestGetTask(t *testing.T) {
	t.Parallel()
	is := is.NewRelaxed(t)

	srv := testServer(t, `{"id":"10001","key":"PROJ-123","fields":{"summary":"Test issue"}}`)
	defer srv.Close()

	client := NewClient("username", "token", srv.URL)

	task, err := client.GetTask("PROJ-123")

	is.NoErr(err)
	is.Equal(task.ID, "PROJ-123")
	is.Equal(task.Description, "Test issue")
}

func TestListTasks(t *testing.T) {
	t.Parallel()
	is := is.NewRelaxed(t)

	srv := testServer(
		t,
		`{"issues":[{"id":"10001","key":"PROJ-123","fields":{"summary":"Test issue"}}]}`,
	)
	defer srv.Close()

	client := NewClient("username", "token", srv.URL)

	tasks, err := client.ListTasks("")

	is.NoErr(err)
	is.Equal(len(tasks), 1)
	is.Equal(tasks[0].ID, "PROJ-123")
	is.Equal(tasks[0].Description, "Test issue")
}

