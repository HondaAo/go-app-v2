package usecase

import (
	"context"

	"github.com/HondaAo/video-app/pkg/auth/model"
)

// Auth repository interface
type UseCase interface {
	Register(ctx context.Context, user *model.User) (*model.UserWithToken, error)
}
