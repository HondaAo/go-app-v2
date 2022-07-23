package usecase

import (
	"context"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/opentracing/opentracing-go"
)

type videoUsecase struct {
	cfg             config.Config
	logger          log.Logger
	videoRepository Repository
}

type Repository interface {
	PostVideo(ctx context.Context, video *model.Video, script []*model.Script) (*model.Video, error)
}

func NewVideoUsecase(cfg config.Config, logger log.Logger, videoRepository Repository) *videoUsecase {
	return &videoUsecase{
		cfg:             cfg,
		logger:          logger,
		videoRepository: videoRepository,
	}
}

func (u videoUsecase) Post(ctx context.Context, video *model.Video, script []*model.Script) (*model.Video, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "video.Create")
	defer span.Finish()

	// user, err := utils.GetUserFromCtx(ctx)
	// if err != nil {
	// 	return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "newsUC.Create.GetUserFromCtx"))
	// }

	createdVideo, err := u.videoRepository.PostVideo(ctx, video, script)
	if err != nil {
		return nil, err
	}

	return createdVideo, nil
}
