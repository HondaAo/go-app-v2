package server

import (
	"context"
	"database/sql"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo        *echo.Echo
	db          *sql.DB
	redisClient *redis.Client
	conf        config.Config
	log         log.Logger
}

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

func NewServer(db *sql.DB, conf config.Config, redisClient *redis.Client, logger log.Logger) *Server {
	return &Server{echo: echo.New(), db: db, redisClient: redisClient, conf: conf, log: logger}
}

func (s Server) Run() error {
	server := &http.Server{
		Addr:           s.conf.Server.Port,
		ReadTimeout:    time.Second * s.conf.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.conf.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.conf.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.log.Fatalf("Error starting Server: ", err)
		}
	}()

	go func() {
		s.log.Infof("Starting Debug Server on PORT: %s", s.conf.Server.PprofPort)
		if err := http.ListenAndServe(s.conf.Server.PprofPort, http.DefaultServeMux); err != nil {
			s.log.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()
	return s.echo.Server.Shutdown(ctx)
}
