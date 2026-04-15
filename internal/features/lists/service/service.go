package lists_service

import (
	"context"

	"github.com/med0viy/practika/internal/core/domain"
)

type ListsService struct {
	listsRepository ListsRepository
}

type ListsRepository interface {
	CreateList(ctx context.Context, list domain.List) (domain.List, error)
	GetLists(ctx context.Context, userID *int) ([]domain.List, error)
	GetList(ctx context.Context, id int) (domain.List, error)
	DeleteList(ctx context.Context, id int) error
	PatchList(ctx context.Context, id int, list domain.List) (domain.List, error)
}

func NewListsService(listsRepository ListsRepository) *ListsService {
	return &ListsService{
		listsRepository: listsRepository,
	}
}
