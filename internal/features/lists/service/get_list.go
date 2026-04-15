package lists_service

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (s *ListsService) GetList(
	ctx context.Context,
	id int,
) (domain.List, error) {
	list, err := s.listsRepository.GetList(ctx, id)
	if err != nil {
		return domain.List{}, fmt.Errorf("get list from repository: %w", err)
	}

	return list, nil
}
