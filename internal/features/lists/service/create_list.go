package lists_service

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (s *ListsService) CreateList(
	ctx context.Context,
	list domain.List,
) (domain.List, error) {
	if err := list.Validate(); err != nil {
		return domain.List{}, fmt.Errorf("validate list domain: %w", err)
	}

	list, err := s.listsRepository.CreateList(ctx, list)
	if err != nil {
		return domain.List{}, fmt.Errorf("create list from repository: %w", err)
	}

	return list, nil
}
