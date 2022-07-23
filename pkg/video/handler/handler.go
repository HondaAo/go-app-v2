package handler

import (
	"encoding/json"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/HondaAo/video-app/pkg/video/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type Handler interface {
	POST() echo.HandlerFunc
}

type videoHandlers struct {
	cfg          *config.Config
	logger       log.Logger
	videoUsecase usecase.Usecase
}

type VideoResponse struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	View     int    `json:"view"`
	Category string `json:"category"`
	Series   string `json:"series"`
	End      int    `json:"end"`
	Start    int    `json:"start" gorm:"default:0"`
	Level    string `json:"level"`
	Script   string `json:"script"`
}

func NewVideoHandler(cfg *config.Config, logger log.Logger, videoUsecase usecase.Usecase) *videoHandlers {
	return &videoHandlers{
		cfg:          cfg,
		logger:       logger,
		videoUsecase: videoUsecase,
	}
}

func (h videoHandlers) POST() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.Create")
		defer span.Finish()

		videoResponse := &VideoResponse{}
		if err := c.Bind(videoResponse); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		var scripts []*model.Script
		if err := json.Unmarshal([]byte(videoResponse.Script), &scripts); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}
		video := &model.Video{
			Title:    videoResponse.Title,
			Url:      videoResponse.Url,
			Category: videoResponse.Category,
			Series:   videoResponse.Series,
			View:     videoResponse.View,
			Start:    videoResponse.Start,
			End:      videoResponse.End,
			Level:    videoResponse.Level,
		}

		createdVideo, err := h.videoUsecase.Post(ctx, video, scripts)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		return c.JSON(200, createdVideo)
	}
}
