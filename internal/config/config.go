package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type secret string

func (s *secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var plainText string
	if err := unmarshal(&plainText); err == nil {
		*s = secret(plainText)
		return nil
	}

	var nested map[string]string
	if err := unmarshal(&nested); err == nil {
		if pass, ok := nested["pass"]; ok {
			p, err := getSecret(pass)
			if err != nil {
				return err
			}
			*s = secret(p)
		}
		return nil
	}

	return fmt.Errorf("failed to unmarshal secret")
}

func (s *secret) String() string {
	return string(*s)
}

func expandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, path[1:]), nil
	}
	return path, nil
}

type Config struct {
	GithubToken secret `yaml:"github_token"`
	TogglToken  secret `yaml:"toggl_token"`
}

func Load(path string) (*Config, error) {
	path, err := expandPath(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
