package usecase

import (
	"context"
	"errors"

	"github.com/HondaAo/video-app/pkg/auth/model"
)

func NewUserModel() *model.User {
	user := &model.User{
		UserID:   "1",
		Email:    "test@test.com",
		Password: "Test",
	}
	user.HashPassword()
	return user
}

type AuthRepositoryMock struct{}

func (a *AuthRepositoryMock) Register(ctx context.Context, user *model.User) (*model.User, error) {
	return NewUserModel(), nil
}

func (a *AuthRepositoryMock) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return NewUserModel(), nil
}

func (a *AuthRepositoryMock) FindByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	return NewUserModel(), nil
}

type AuthRepositoryMockError struct{}

func (a *AuthRepositoryMockError) Register(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, errors.New("register Error")
}

func (a *AuthRepositoryMockError) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return nil, errors.New("get Error")
}

func (a *AuthRepositoryMockError) FindByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, errors.New("get Error")
}

type RedisRepositoryMock struct{}

func (r *RedisRepositoryMock) GetByIDCtx(ctx context.Context, key string) (*model.User, error) {
	return NewUserModel(), nil
}

func (r *RedisRepositoryMock) SetUserCtx(ctx context.Context, key string, seconds int, user *model.User) error {
	return nil
}

type RedisRepositoryMockNil struct{}

func (r *RedisRepositoryMockNil) GetByIDCtx(ctx context.Context, key string) (*model.User, error) {
	return nil, nil
}

func (r *RedisRepositoryMockNil) SetUserCtx(ctx context.Context, key string, seconds int, user *model.User) error {
	return nil
}

type RedisRepositoryMockError struct{}

func (r *RedisRepositoryMockError) GetByIDCtx(ctx context.Context, key string) (*model.User, error) {
	return nil, errors.New("get error")
}

func (r *RedisRepositoryMockError) SetUserCtx(ctx context.Context, key string, seconds int, user *model.User) error {
	return errors.New("get error")
}
