package handler

import "github.com/labstack/echo/v4"

type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	GetUsers() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	GetCSRFToken() echo.HandlerFunc
}

func MapAuthRoutes(authGroup *echo.Group, authHandler Handlers) {
	authGroup.POST("/register", authHandler.Register())
}
