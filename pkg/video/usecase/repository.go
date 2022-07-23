package usecase

import (
	"context"

	"github.com/HondaAo/video-app/pkg/video/model"
)

type Usecase interface {
	Post(ctx context.Context, video *model.Video, script []*model.Script) (*model.Video, error)
}
