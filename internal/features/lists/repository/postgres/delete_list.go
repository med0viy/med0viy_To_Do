package lists_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/med0viy/practika/internal/core/errors"
)

func (r *ListsRepository) DeleteList(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.lists
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"list with id='%d': %w",
			id,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
