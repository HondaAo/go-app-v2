package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/HondaAo/video-app/config"
	logger "github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/server"
	"github.com/HondaAo/video-app/utils"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
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

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	appLogger.Info("Opentracing connected")

	explicit("digital-maker-270904")

	s := server.NewServer(mysql, *cfg, redisClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}

}

func explicit(projectID string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/app/cmd/api/google.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("Buckets:")
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(battrs.Name)
	}
}
