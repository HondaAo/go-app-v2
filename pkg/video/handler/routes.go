package handler

import (
	"github.com/HondaAo/video-app/pkg/auth/middleware"
	"github.com/labstack/echo/v4"
)

func MapVideoRoute(videoGroup echo.Group, h Handler, mw middleware.MiddlewareManager) {
	videoGroup.POST("/create", h.POST(), mw.AuthSessionMiddleware)
}
