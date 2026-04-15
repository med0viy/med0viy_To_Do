package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	taskID int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET
		version=version + 1,
		title=$1,
		description=$2,
		complited=$3,
		is_important=$4,
		is_in_my_day=$5,
		due_date=$6,
		complited_at=$7,
		list_id=$8
	WHERE id=$9 AND version=$10
	RETURNING 
		id,
		version,
		title,
		description,
		complited,
		is_important,
		is_in_my_day,
		created_at,
		due_date,
		complited_at,
		list_id,
		author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Complited,
		task.IsImportant,
		task.IsInMyDay,
		task.DueDate,
		task.ComplitedAt,
		task.ListID,
		taskID,
		task.Version,
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
		&taskModel.ListID,
		&taskModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d' concurrently accessed: %w",
				taskID,
				core_errors.ErrConflict,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := TaskDomainFromModel(taskModel)

	return taskDomain, nil
}
