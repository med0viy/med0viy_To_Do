package domain

import (
	"fmt"
	"time"

	core_errors "github.com/med0viy/practika/internal/core/errors"
)

type Task struct {
	ID          int
	Version     int
	Title       string
	Description *string
	Complited   bool
	IsImportant bool
	IsInMyDay   bool
	CreatedAt   time.Time
	DueDate     *time.Time
	ComplitedAt *time.Time
	AutorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	complited bool,
	isImportant bool,
	isInMyDay bool,
	createdAt time.Time,
	dueDate *time.Time,
	complitedAt *time.Time,
	autorUserID int,
) Task {
	return Task{
		ID:          id,
		Version:     version,
		Title:       title,
		Description: description,
		Complited:   complited,
		IsImportant: isImportant,
		IsInMyDay:   isInMyDay,
		CreatedAt:   createdAt,
		DueDate:     dueDate,
		ComplitedAt: complitedAt,
		AutorUserID: autorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	isImportant bool,
	isInMyDay bool,
	dueDate *time.Time,
	autorUserID int,
) Task {
	return NewTask(
		UninitiolizedID,
		UninitiolizedVersion,
		title,
		description,
		false,
		isImportant,
		isInMyDay,
		time.Now(),
		dueDate,
		nil,
		autorUserID,
	)
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Complited {
		if t.ComplitedAt == nil {
			return fmt.Errorf(
				"`ComplitedAt` can't be `nil` if `Complited`==`true`: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.ComplitedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`ComplitedAt` can't be before `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.ComplitedAt != nil {
			return fmt.Errorf(
				"`ComplitedAt` must be `nil` if `Complited`==`false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Complited   Nullable[bool]
	IsImportant Nullable[bool]
	IsInMyDay   Nullable[bool]
	DueDate     Nullable[time.Time]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	complited Nullable[bool],
	isImportant Nullable[bool],
	isInMyDay Nullable[bool],
	dueDate Nullable[time.Time],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Complited:   complited,
		IsImportant: isImportant,
		IsInMyDay:   isInMyDay,
		DueDate:     dueDate,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"`Title` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Complited.Set && p.Complited.Value == nil {
		return fmt.Errorf(
			"`Complited` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.IsImportant.Set && p.IsImportant.Value == nil {
		return fmt.Errorf(
			"`IsImportant` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.IsInMyDay.Set && p.IsInMyDay.Value == nil {
		return fmt.Errorf(
			"`IsInMyDay` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (t *Task) ApplyPatched(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Complited.Set {
		tmp.Complited = *patch.Complited.Value

		if tmp.Complited {
			complitedAt := time.Now()
			tmp.ComplitedAt = &complitedAt
		} else {
			tmp.ComplitedAt = nil
		}
	}

	if patch.IsImportant.Set {
		tmp.IsImportant = *patch.IsImportant.Value
	}

	if patch.IsInMyDay.Set {
		tmp.IsInMyDay = *patch.IsInMyDay.Value
	}

	if patch.DueDate.Set {
		tmp.DueDate = patch.DueDate.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate task after patch apply: %w", err)
	}

	*t = tmp

	return nil
}
