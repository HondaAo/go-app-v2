package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/HondaAo/video-app/pkg/video/usecase"
	"github.com/HondaAo/video-app/utils"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type Handler interface {
	Post() echo.HandlerFunc
	Get() echo.HandlerFunc
	GetAll() echo.HandlerFunc
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

func (h videoHandlers) Post() echo.HandlerFunc {
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

func (h videoHandlers) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "videoHandlers.GETALL")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		videoList, err := h.videoUsecase.GetAll(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, videoList)
	}
}

func (h videoHandlers) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "videoHandlers.GET")
		defer span.Finish()

		id, _ := strconv.Atoi(c.Param("id"))

		video, script, err := h.videoUsecase.Get(ctx, id)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(utils.ErrorResponse(err))
		}

		stringScript, _ := json.Marshal(script)
		videoResponse := &VideoResponse{
			Title:    video.Title,
			Url:      video.Url,
			Category: video.Category,
			Series:   video.Series,
			View:     video.View,
			Start:    video.Start,
			End:      video.End,
			Script:   string(stringScript),
		}

		return c.JSON(http.StatusOK, videoResponse)
	}
}
