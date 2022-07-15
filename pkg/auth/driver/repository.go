package driver

import (
	"context"

	"github.com/HondaAo/video-app/pkg/auth/model"
)

type Repository interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	// Update(ctx context.Context, user *models.User) (*models.User, error)
	// Delete(ctx context.Context, userID uuid.UUID) error
	// GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	FindByEmail(ctx context.Context, user *model.User) (*model.User, error)
	// GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}
