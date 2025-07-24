package project

import (
	"regexp"

	"github.com/destag/ttrack/internal/config"
)

type Project struct {
	Name         string
	TaskID       string
	Type         string
	Source       string
	BranchFormat string
}

func Find(projects map[string]config.Project, input string) (Project, bool) {
	for name, pr := range projects {
		for _, task := range pr.Tasks {
			re := regexp.MustCompile(task.Regex)
			if matches := re.FindStringSubmatch(input); len(matches) > 0 {
				return Project{
					Name:         name,
					TaskID:       matches[0],
					Type:         task.Type,
					Source:       task.Source,
					BranchFormat: pr.BranchFormat,
				}, true
			}
		}
	}
	return Project{}, false
}
