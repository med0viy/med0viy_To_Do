package lists_service

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (s *ListsService) GetLists(
	ctx context.Context,
	userID *int,
) ([]domain.List, error) {
	lists, err := s.listsRepository.GetLists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get lists from repository: %w", err)
	}

	return lists, nil
}
