package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description, complited, is_important, is_in_my_day, created_at, due_date, complited_at, autor_user_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, version, title, description, complited, is_important, is_in_my_day, created_at, due_date, complited_at, autor_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Complited,
		task.IsImportant,
		task.IsInMyDay,
		task.CreatedAt,
		task.DueDate,
		task.ComplitedAt,
		task.AutorUserID,
	)

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
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v: user with id='%d': %w",
				err,
				task.AutorUserID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := TaskDomainFromModel(taskModel)

	return taskDomain, nil
}
