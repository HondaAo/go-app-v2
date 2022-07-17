package handler

import (
	dLog "log"
	"net/http"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/auth/model"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	// Logout() echo.HandlerFunc
	// Update() echo.HandlerFunc
	// Delete() echo.HandlerFunc
	// GetUserByID() echo.HandlerFunc
	// GetUsers() echo.HandlerFunc
	// GetMe() echo.HandlerFunc
	// GetCSRFToken() echo.HandlerFunc
}

// Auth handlers
type authHandlers struct {
	cfg    *config.Config
	authUC usecase.UseCase
	logger log.Logger
	redis  redis.Client
}

// NewAuthHandlers Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC usecase.UseCase, log log.Logger, redis redis.Client) Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log, redis: redis}
}

// Register
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Register")
		defer span.Finish()

		user := &model.User{}
		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		createdUser, err := h.authUC.Register(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(500, err)
		}

		sess, err := utils.CreateSession(ctx, h.redis, utils.Session{
			UserID: createdUser.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(500, err)
		}

		c.SetCookie(&http.Cookie{
			Name:       h.cfg.Session.Name,
			Value:      sess,
			Path:       "/",
			RawExpires: "",
			MaxAge:     h.cfg.Session.Expire,
			Secure:     h.cfg.Cookie.Secure,
			HttpOnly:   h.cfg.Cookie.HTTPOnly,
			SameSite:   0,
		})

		return c.JSON(http.StatusCreated, createdUser)
	}
}

func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Register")
		defer span.Finish()

		login := &Login{}
		if err := c.Bind(&login); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		userWithToken, err := h.authUC.Login(ctx, &model.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		sess, err := utils.CreateSession(ctx, h.redis, utils.Session{
			UserID: userWithToken.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			dLog.Print(err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		c.SetCookie(&http.Cookie{
			Name:       h.cfg.Session.Name,
			Value:      sess,
			Path:       "/",
			RawExpires: "",
			MaxAge:     h.cfg.Session.Expire,
			Secure:     h.cfg.Cookie.Secure,
			HttpOnly:   h.cfg.Cookie.HTTPOnly,
			SameSite:   0,
		})

		return c.JSON(http.StatusOK, userWithToken)
	}
}
