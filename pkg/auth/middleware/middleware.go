package middleware

import (
	"context"
	"net/http"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// Middleware manager
type MiddlewareManager struct {
	authUC  usecase.UseCase
	cfg     *config.Config
	origins []string
	logger  log.Logger
	redis   redis.Client
}

// Middleware manager constructor
func NewMiddlewareManager(authUC usecase.UseCase, cfg *config.Config, origins []string, logger log.Logger, redis redis.Client) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, origins: origins, logger: logger, redis: redis}
}

// Auth sessions middleware using redis
func (mw *MiddlewareManager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(mw.cfg.Session.Name)
		if err != nil {
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError(utils.Unauthorized))
		}

		sid := cookie.Value

		sess, err := utils.GetSessionByID(c.Request().Context(), mw.redis, cookie.Value)
		if err != nil {
			mw.logger.Errorf("GetSessionByID RequestID: %s, CookieValue: %s, Error: %s",
				utils.GetRequestID(c),
				cookie.Value,
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError(utils.Unauthorized))
		}

		user, err := mw.authUC.GetByID(c.Request().Context(), sess.UserID)
		if err != nil {
			mw.logger.Errorf("GetByID RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError(utils.Unauthorized))
		}

		c.Set("sid", sid)
		c.Set("uid", sess.SessionID)
		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		mw.logger.Info(
			"SessionMiddleware, RequestID: %s,  IP: %s, UserID: %s, CookieSessionID: %s",
			utils.GetRequestID(c),
			utils.GetIPAddress(c),
			user.UserID,
			cookie.Value,
		)

		return next(c)
	}
}
