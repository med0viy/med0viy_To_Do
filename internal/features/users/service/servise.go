package users_service

import (
	"context"

	"github.com/med0viy/practika/internal/core/domain"
)

type UsersServise struct {
	usersRepository usersRepository
}

type usersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error)
}

func NewUsersServise(usersRepository UsersRepository) *UsersServise {
	return &UsersServise{
		usersRepository: usersRepository,
	}
}
