package lists_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
)

func (r *ListsRepository) GetList(
	ctx context.Context,
	id int,
) (domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, author_user_id
	FROM todoapp.lists
	WHERE id=$1;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		id,
	)

	var listModel ListModel

	err := row.Scan(
		&listModel.ID,
		&listModel.Version,
		&listModel.Name,
		&listModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.List{}, fmt.Errorf(
				"list with id='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return domain.List{}, fmt.Errorf("scan error: %w", err)
	}

	listDomain := ListDomainFromModel(listModel)

	return listDomain, nil
}
