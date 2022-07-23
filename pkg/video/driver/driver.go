package driver

import (
	"context"
	"database/sql"

	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type videoRepo struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *videoRepo {
	return &videoRepo{
		db: db,
	}
}

func (r videoRepo) PostVideo(ctx context.Context, video *model.Video, scripts []*model.Script) (*model.Video, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Create")
	defer span.Finish()

	v := &model.Video{}
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO videos(title,url,category,series,end,start,level) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return nil, errors.Wrap(err, "Video.Insert Error")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, &video.Title, &video.Url, &video.Category, &video.Series, &video.End, &video.Start, &video.Level)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()

	for _, script := range scripts {
		stmt, err := r.db.PrepareContext(ctx, "INSERT INTO scripts(video_id,text,ja,timestamp) VALUES(?,?,?,?)")
		if err != nil {
			return nil, errors.Wrap(err, "Script.Insert Error")
		}
		defer stmt.Close()

		res, _ = stmt.ExecContext(ctx, lastId, &script.Text, &script.Ja, &script.TimeStamp)
	}
	return v, nil
}
