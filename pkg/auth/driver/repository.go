package driver

import (
	"context"

	"github.com/HondaAo/video-app/pkg/auth/model"
)

type Repository interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	// Update(ctx context.Context, user *models.User) (*models.User, error)
	// Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID string) (*model.User, error)
	// FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	FindByEmail(ctx context.Context, user *model.User) (*model.User, error)
	// GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*model.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *model.User) error
	// DeleteUserCtx(ctx context.Context, key string) error
}
