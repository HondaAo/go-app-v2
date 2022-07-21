package server

import (
	"log"

	"github.com/HondaAo/video-app/pkg/auth/driver"
	"github.com/HondaAo/video-app/pkg/auth/handler"
	authHandler "github.com/HondaAo/video-app/pkg/auth/handler"
	apiMiddleware "github.com/HondaAo/video-app/pkg/auth/middleware"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MapHandler(e *echo.Echo) error {

	aRepo := driver.NewAuthRepository(s.db)
	authRedisRepo := driver.NewAuthRedisRepo(s.redisClient)

	aUsecase := usecase.NewAuthUseCase(&s.conf, aRepo, authRedisRepo, s.log)

	log.Print(s.redisClient)
	aHandler := handler.NewAuthHandlers(&s.conf, aUsecase, s.log, *s.redisClient)

	mw := apiMiddleware.NewMiddlewareManager(aUsecase, &s.conf, []string{"*"}, s.log, *s.redisClient)

	e.Use(mw.RequestLoggerMiddleware)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, utils.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	v1 := e.Group("/api/v1")
	authGroup := v1.Group("/auth")

	authHandler.MapAuthRoutes(authGroup, aHandler, *mw)

	return nil
}
