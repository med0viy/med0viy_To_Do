package tasks_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/med0viy/practika/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	listID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id, version, title, description, complited, is_important, is_in_my_day, created_at, due_date, complited_at, list_id, author_user_id
	FROM todoapp.tasks
	`)

	var conditions []string
	var args []any
	argID := 1

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", argID))
		args = append(args, userID)
		argID++
	}

	if listID != nil {
		conditions = append(conditions, fmt.Sprintf("list_id = $%d", argID))
		args = append(args, listID)
		argID++
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC")
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID + 1))
	args = append(args, limit, offset)

	rows, err := r.pool.Query(
		ctx,
		queryBuilder.String(),
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
			&taskModel.ListID,
			&taskModel.AuthorUserID,
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
