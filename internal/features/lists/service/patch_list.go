package lists_service

import (
	"context"
	"fmt"

	"github.com/med0viy/practika/internal/core/domain"
)

func (s *ListsService) PatchList(
	ctx context.Context,
	id int,
	patch domain.ListPatch,
) (domain.List, error) {
	list, err := s.listsRepository.GetList(ctx, id)
	if err != nil {
		return domain.List{}, fmt.Errorf(
			"get list from repository: %w",
			err,
		)
	}

	if err := list.ApplyPatch(patch); err != nil {
		return domain.List{}, fmt.Errorf(
			"apply patch error: %w",
			err,
		)
	}

	patchedList, err := s.listsRepository.PatchList(ctx, id, list)
	if err != nil {
		return domain.List{}, fmt.Errorf(
			"patch list from repository: %w",
			err,
		)
	}

	return patchedList, nil
}
