package lists_postgres_repository

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (r *ListsRepository) GetLists(
	ctx context.Context,
	userID *int,
) ([]domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, author_user_id
	FROM todoapp.lists
	%s
	ORDER BY id ASC;
	`
	args := []any{}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$1")
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
		return nil, fmt.Errorf("select lists: %w", err)	
	}
	defer rows.Close()

	listsModels := []ListModel{}

	for rows.Next() {
		var listModel ListModel

		err := rows.Scan(
			&listModel.ID,
			&listModel.Version,
			&listModel.Name,
			&listModel.AuthorUserID,
		)

		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		listsModels = append(listsModels, listModel)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}
	
	listsDomain := ListsDomainsFromModels(listsModels)

	return listsDomain, nil
}
