package tasks_postgres_repository

import (
	"time"

	"github.com/med0viy/practika/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Complited    bool
	IsImportant  bool
	IsInMyDay    bool
	CreatedAt    time.Time
	DueDate      *time.Time
	ComplitedAt  *time.Time
	ListID       *int
	AuthorUserID int
}

func TaskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Complited,
		taskModel.IsImportant,
		taskModel.IsInMyDay,
		taskModel.CreatedAt,
		taskModel.DueDate,
		taskModel.ComplitedAt,
		taskModel.ListID,
		taskModel.AuthorUserID,
	)
}

func TaskDomainsFromModels(tasks []TaskModel) []domain.Task {
	taskDomains := make([]domain.Task, len(tasks))

	for k, model := range tasks {
		taskDomains[k] = TaskDomainFromModel(model)
	}

	return taskDomains
}
