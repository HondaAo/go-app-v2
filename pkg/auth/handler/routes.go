package handler

import (
	"github.com/HondaAo/video-app/pkg/auth/middleware"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, authHandler Handlers, mw middleware.MiddlewareManager) {
	authGroup.POST("/register", authHandler.Register())
	authGroup.POST("/login", authHandler.Login())
	authGroup.POST("/logout", authHandler.Logout())
	authGroup.Use(mw.AuthSessionMiddleware)
	authGroup.GET("/me", authHandler.GetMe())
}
