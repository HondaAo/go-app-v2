package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/auth/model"
	authUtils "github.com/HondaAo/video-app/pkg/auth/utils"
	"github.com/HondaAo/video-app/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg       *config.Config
	authRepo  Repository
	redisRepo RedisRepository
	logger    log.Logger
}

type Repository interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*model.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *model.User) error
	// DeleteUserCtx(ctx context.Context, key string) error
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo Repository, redisRepo RedisRepository, log log.Logger) UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, logger: log}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *model.User) (*model.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser.FirstName != "" || err != nil {
		return nil, utils.NewRestErrorWithMessage(http.StatusBadRequest, utils.ErrEmailAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, utils.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := authUtils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, utils.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	createdUser.Password = ""

	return &model.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *model.User) (*model.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	defer span.Finish()

	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, utils.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
	}

	token, err := authUtils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, utils.NewInternalServerError(errors.Wrap(err, "authUC.GetUsers.GenerateJWTToken"))
	}

	foundUser.Password = ""

	return &model.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, userID string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	defer span.Finish()

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID))
	if err != nil {
		// u.logger.Errorf("authUC.GetByID.GetByIDCtx: %v", err)
		return nil, err
	}
	if cachedUser != nil {
		cachedUser.Password = ""
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
	}

	user.Password = ""

	return user, nil
}

func (u *authUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
