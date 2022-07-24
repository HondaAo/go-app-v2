package usecase

import (
	"context"

	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/HondaAo/video-app/utils"
)

type Usecase interface {
	Post(ctx context.Context, video *model.Video, script []*model.Script) (*model.Video, error)
	GetAll(ctx context.Context, pq *utils.PaginationQuery) ([]*model.Video, error)
	Get(ctx context.Context, id int) (*model.Video, []*model.Script, error)
}
