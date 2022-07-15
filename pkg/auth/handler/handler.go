package handler

import (
	"net/http"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/auth/model"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type Handlers interface {
	Register() echo.HandlerFunc
	// Login() echo.HandlerFunc
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
}

// NewAuthHandlers Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC usecase.UseCase, log log.Logger) Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log}
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

		sess, err := utils.CreateSession(ctx, utils.Session{
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
