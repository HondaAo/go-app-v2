package main

import (
	"log"
	"os"

	"github.com/HondaAo/video-app/config"
	logger "github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/server"
	"github.com/HondaAo/video-app/utils"
)

func main() {
	configPath := config.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	mysql := utils.NewDB()
	defer mysql.Close()
	appLogger.Info("Mysql connected")

	redisClient := utils.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	s := server.NewServer(mysql, *cfg, redisClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
