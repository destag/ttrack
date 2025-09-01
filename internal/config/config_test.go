package config

import (
	"testing"

	"github.com/matryer/is"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	is := is.NewRelaxed(t)

	cfg, err := Load("testdata/config.yml")

	is.NoErr(err)
	is.Equal(cfg.GithubToken.String(), "ghp_secret")
	is.Equal(cfg.TogglToken.String(), "toggl_secret")
}
