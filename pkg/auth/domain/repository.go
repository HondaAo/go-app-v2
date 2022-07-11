package domain

import (
	"context"

	"github.com/HondaAo/video-app/pkg/auth/domain/entity"
	"github.com/HondaAo/video-app/pkg/auth/driver/model"
)

// Auth repository interface
type Repository interface {
	Register(ctx context.Context, user *model.User) (*entity.UserEntityWithToken, error)
	// Login(ctx context.Context, user *entity.UserEntity) (*entity.UserEntityWithToken, error)
	// Update(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error)
	// Delete(ctx context.Context, userID uuid.UUID) error
	// GetByID(ctx context.Context, userID uuid.UUID) (*entity.UserEntity, error)
}
