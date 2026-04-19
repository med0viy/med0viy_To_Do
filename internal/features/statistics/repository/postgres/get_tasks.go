package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/med0viy/practika/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userID *int,
	listID *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error){
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id, version, title, description, complited, is_important, is_in_my_day, created_at, due_date, complited_at, list_id, author_user_id
	FROM todoapp.tasks
	`)

	var conditions []string
	var args []any
	argsID := 1

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", argsID))
		args = append(args, userID)
		argsID++
	}

	if listID != nil {
		conditions = append(conditions, fmt.Sprintf("list_id = $%d", argsID))
		args = append(args, listID)
		argsID++
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argsID))
		args = append(args, from)
		argsID++
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", argsID))
		args = append(args, to)
		argsID++
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var tasksModels []TaskModel
	
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
			return nil, fmt.Errorf("scan error: %w", err)
		}

		tasksModels = append(tasksModels, taskModel)
	}

	if  err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	tasksDomains := TaskDomainsFromModels(tasksModels)

	return tasksDomains, nil
}
