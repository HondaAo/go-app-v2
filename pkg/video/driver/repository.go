package driver

import (
	"context"

	"github.com/HondaAo/video-app/pkg/video/model"
)

type Repository interface {
	PostVideo(ctx context.Context, video *model.Video, script *model.Script) (*model.Video, error)
}
