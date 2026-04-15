package lists_transport_http

import (
	"context"
	"net/http"

	"github.com/med0viy/practika/internal/core/domain"
	core_http_server "github.com/med0viy/practika/internal/core/transport/http/server"
)

type ListsHTTPHandler struct {
	listsService ListsService
}

type ListsService interface {
	CreateList(ctx context.Context, list domain.List) (domain.List, error)
	GetLists(ctx context.Context, userID *int) ([]domain.List, error)
	GetList(ctx context.Context, id int) (domain.List, error)
	DeleteList(ctx context.Context, id int) error
	PatchList(ctx context.Context, id int, patch domain.ListPatch) (domain.List, error)
}

func NewListsHTTPHandler(listsService ListsService) *ListsHTTPHandler {
	return &ListsHTTPHandler{
		listsService: listsService,
	}
}

func (h *ListsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/lists",
			Handler: h.CreateList,
		},
		{
			Method:  http.MethodGet,
			Path:    "/lists",
			Handler: h.GetLists,
		},
		{
			Method:  http.MethodGet,
			Path:    "/lists/{id}",
			Handler: h.GetList,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/lists/{id}",
			Handler: h.DeleteList,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/lists/{id}",
			Handler: h.PatchList,
		},
	}
}
