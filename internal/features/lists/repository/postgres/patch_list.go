package lists_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
)

func (r *ListsRepository) PatchList(
	ctx context.Context,
	id int,
	list domain.List,
) (domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.lists
	SET
		name=$1
	WHERE id=$2 AND version=$3
	RETURNING id, version, name, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		list.Name,
		id,
		list.Version,
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
				"list with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}

		return domain.List{}, fmt.Errorf("scan error: %w", err)
	}

	listDomain := ListDomainFromModel(listModel)

	return listDomain, nil
}
