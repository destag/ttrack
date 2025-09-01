package tasks

type Task struct {
	ID          string
	Description string
}

type TaskProvider interface {
	GetTask(id string) (*Task, error)
	ListTasks(query string) ([]*Task, error)
}
