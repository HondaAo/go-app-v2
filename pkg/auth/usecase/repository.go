package usecase

import (
	"context"

	"github.com/HondaAo/video-app/pkg/auth/model"
)

// Auth repository interface
type UseCase interface {
	Register(ctx context.Context, user *model.User) (*model.UserWithToken, error)
	Login(ctx context.Context, user *model.User) (*model.UserWithToken, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}
