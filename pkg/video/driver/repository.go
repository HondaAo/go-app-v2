package driver

import (
	"context"

	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/HondaAo/video-app/utils"
)

type Repository interface {
	PostVideo(ctx context.Context, video *model.Video, script *model.Script) (*model.Video, error)
	GetVideos(ctx context.Context, pq utils.PaginationQuery) ([]*model.Video, error)
	GetVideo(ctx context.Context, id int) (*model.Video, []*model.Script, error)
}
