package middleware

import (
	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/auth/usecase"
)

// Middleware manager
type MiddlewareManager struct {
	authUC  usecase.UseCase
	cfg     *config.Config
	origins []string
	logger  log.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(authUC usecase.UseCase, cfg *config.Config, origins []string, logger log.Logger) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
