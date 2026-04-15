package lists_transport_http

import "github.com/med0viy/practika/internal/core/domain"

type ListDTOResponse struct {
	ID           int    `json:"id"`
	Version      int    `json:"version"`
	Name         string `json:"name"`
	AuthorUserID int    `json:"author_user_id"`
}

func listDTOFromDomain(list domain.List) ListDTOResponse {
	return ListDTOResponse{
		ID:           list.ID,
		Version:      list.Version,
		Name:         list.Name,
		AuthorUserID: list.AuthorUserID,
	}
}

func listsDTOFromDmains(lists []domain.List) []ListDTOResponse {
	listsDTO := make([]ListDTOResponse, len(lists))

	for k, v := range lists {
		listsDTO[k] = listDTOFromDomain(v)
	}

	return listsDTO
}
