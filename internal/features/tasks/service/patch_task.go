package tasks_service

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	taskID int,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatched(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task from repository: %w", err)
	}

	return patchedTask, nil
}
