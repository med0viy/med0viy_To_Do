package domain

import (
	"fmt"

	core_errors "github.com/med0viy/practika/internal/core/errors"
)

type List struct {
	ID           int
	Version      int
	Name         string
	AuthorUserID int
}

func NewList(
	id int,
	version int,
	name string,
	authorUserID int,
) List {
	return List{
		ID:           id,
		Version:      version,
		Name:         name,
		AuthorUserID: authorUserID,
	}
}

func NewUninitiolizedList(
	name string,
	authorUserID int,
) List {
	return List{
		ID:           UninitiolizedID,
		Version:      UninitiolizedVersion,
		Name:         name,
		AuthorUserID: authorUserID,
	}
}

func (l *List) Validate() error {
	nameLen := len([]rune(l.Name))
	if nameLen < 1 || nameLen > 100 {
		return fmt.Errorf(
			"invalid `Name` len: %d: %w",
			nameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type ListPatch struct {
	Name Nullable[string]
}

func NewListPatch(name Nullable[string]) ListPatch {
	return ListPatch{
		Name: name,
	}
}

func (p *ListPatch) Validate() error {
	if p.Name.Set {
		if p.Name.Value == nil {
			return fmt.Errorf(
				"`Name` can't be patched to NULL: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}


	return nil
}

func (l *List) ApplyPatch(patch ListPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate patch: %w", err)
	}

	tmp := *l

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched list: %w", err)
	}

	*l = tmp

	return nil
}
