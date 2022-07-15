package handler

import "github.com/labstack/echo/v4"

func MapAuthRoutes(authGroup *echo.Group, authHandler Handlers) {
	authGroup.POST("/register", authHandler.Register())
}
