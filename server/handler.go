package server

import (
	"log"

	"github.com/HondaAo/video-app/pkg/auth/driver"
	"github.com/HondaAo/video-app/pkg/auth/handler"
	authHandler "github.com/HondaAo/video-app/pkg/auth/handler"
	apiMiddleware "github.com/HondaAo/video-app/pkg/auth/middleware"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
	videoDriver "github.com/HondaAo/video-app/pkg/video/driver"
	videoHandler "github.com/HondaAo/video-app/pkg/video/handler"
	videoUsecase "github.com/HondaAo/video-app/pkg/video/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MapHandler(e *echo.Echo) error {

	aRepo := driver.NewAuthRepository(s.db)
	vRepo := videoDriver.NewVideoRepository(s.db)
	authRedisRepo := driver.NewAuthRedisRepo(s.redisClient)

	aUsecase := usecase.NewAuthUseCase(&s.conf, aRepo, authRedisRepo, s.log)
	vUsecase := videoUsecase.NewVideoUsecase(s.conf, s.log, vRepo)

	log.Print(s.redisClient)
	aHandler := handler.NewAuthHandlers(&s.conf, aUsecase, s.log, *s.redisClient)
	vHandler := videoHandler.NewVideoHandler(&s.conf, s.log, vUsecase)

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
	e.Use(middleware.Secure())

	v1 := e.Group("/api/v1")
	authGroup := v1.Group("/auth")
	videoGroup := v1.Group("/video")

	authHandler.MapAuthRoutes(authGroup, aHandler, *mw)
	videoHandler.MapVideoRoute(*videoGroup, vHandler, *mw)

	return nil
}
