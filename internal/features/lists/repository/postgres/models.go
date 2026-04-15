package lists_postgres_repository

import "github.com/med0viy/practika/internal/core/domain"

type ListModel struct {
	ID           int
	Version      int
	Name         string
	AuthorUserID int
}

func ListDomainFromModel(list ListModel) domain.List {
	return domain.NewList(
		list.ID,
		list.Version,
		list.Name,
		list.AuthorUserID,
	)
}

func ListsDomainsFromModels(lists []ListModel) []domain.List {
	listsDomains := make([]domain.List, len(lists))

	for k, model := range lists {
		listsDomains[k] = ListDomainFromModel(model)
	}

	return listsDomains
}
