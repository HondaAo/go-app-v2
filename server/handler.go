package server

import (
	"github.com/HondaAo/video-app/pkg/auth/driver"
	"github.com/HondaAo/video-app/pkg/auth/handler"
	authHandler "github.com/HondaAo/video-app/pkg/auth/handler"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandler(e *echo.Echo) error {

	aRepo := driver.NewAuthRepository(s.db)

	aUsecase := usecase.NewAuthUseCase(&s.conf, aRepo, s.log)

	aHandler := handler.NewAuthHandlers(&s.conf, aUsecase, s.log)

	v1 := e.Group("/api/v1")
	authGroup := v1.Group("/auth")

	authHandler.MapAuthRoutes(authGroup, aHandler)

	return nil
}
