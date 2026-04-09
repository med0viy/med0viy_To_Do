package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, complited, is_important, is_in_my_day, created_at, due_date, complited_at, autor_user_id
	FROM todoapp.tasks
	%s
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	args := []any{limit, offset}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE autor_user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf("select tasks:  %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := TaskDomainsFromModels(taskModels)

	return taskDomains, nil
}
