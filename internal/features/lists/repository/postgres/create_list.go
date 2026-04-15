package lists_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
)

func (r *ListsRepository) CreateList(
	ctx context.Context,
	list domain.List,
) (domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.lists (name, author_user_id)
	VALUES($1, $2)
	RETURNING id, version, name, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		list.Name,
		list.AuthorUserID,
	)

	var listModel ListModel
	err := row.Scan(
		&listModel.ID,
		&listModel.Version,
		&listModel.Name,
		&listModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.List{}, fmt.Errorf(
				"%v: user with id='%d': %w",
				err,
				list.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}

		return domain.List{}, fmt.Errorf("scan error: %w", err)
	}

	listDomain := ListDomainFromModel(listModel)

	return listDomain, nil
}
