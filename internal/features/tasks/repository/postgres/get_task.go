package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, complited, is_important, is_in_my_day, created_at, due_date, complited_at, autor_user_id
	FROM todoapp.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, taskID)

	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Complited,
		&taskModel.IsImportant,
		&taskModel.IsInMyDay,
		&taskModel.CreatedAt,
		&taskModel.DueDate,
		&taskModel.ComplitedAt,
		&taskModel.AutorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d': %w",
				taskID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := TaskDomainFromModel(taskModel)

	return taskDomain, nil
}
