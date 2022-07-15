package usecase

import (
	"context"
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
	cfg      *config.Config
	authRepo Repository
	logger   log.Logger
}

type Repository interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, user *model.User) (*model.User, error)
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo Repository, log log.Logger) UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, logger: log}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *model.User) (*model.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, utils.NewRestErrorWithMessage(http.StatusBadRequest, utils.ErrEmailAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, utils.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.Password = ""

	token, err := authUtils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, utils.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &model.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}
