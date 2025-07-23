package project

import (
	"regexp"

	"github.com/destag/ttrack/internal/config"
)

func Find(projects map[string]config.Project, input string) (config.Project, string, bool) {
	for rgx, pr := range projects {
		re := regexp.MustCompile(rgx)
		if matches := re.FindStringSubmatch(input); len(matches) > 0 {
			return pr, matches[0], true
		}
	}
	return config.Project{}, "", false
}
