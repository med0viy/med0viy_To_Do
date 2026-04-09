package tasks_service

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from reposiroty: %w", err)
	}

	return task, nil
}
