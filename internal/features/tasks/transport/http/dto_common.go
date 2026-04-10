package tasks_transport_http

import (
	"time"

	"github.com/med0viy/practika/internal/core/domain"
)

type TaskDTOResponse struct {
	Id           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Complited    bool       `json:"complited"`
	IsImportant  bool       `json:"is_important"`
	IsInMyDay    bool       `json:"is_in_my_day"`
	CreatedAt    time.Time  `json:"created_at"`
	DueDate      *time.Time `json:"due_date"`
	ComplitedAt  *time.Time `json:"complited_at"`
	AuthorUserID int        `json:"author_user_id"`
}

func TaskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		Id:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Complited:    task.Complited,
		IsImportant:  task.IsImportant,
		IsInMyDay:    task.IsInMyDay,
		CreatedAt:    task.CreatedAt,
		DueDate:      task.DueDate,
		ComplitedAt:  task.ComplitedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func TasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))

	for k, v := range tasks {
		tasksDTO[k] = TaskDTOFromDomain(v)
	}

	return tasksDTO
}
