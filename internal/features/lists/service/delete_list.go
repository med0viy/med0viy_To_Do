package lists_service

import (
	"context"
	"fmt"
)

func (s *ListsService) DeleteList(
	ctx context.Context,
	id int,
) error {
	if err := s.listsRepository.DeleteList(ctx, id); err != nil {
		return fmt.Errorf("delete list from repository: %w", err)
	}

	return nil
}
