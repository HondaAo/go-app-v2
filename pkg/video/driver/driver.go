package driver

import (
	"context"
	"database/sql"

	"github.com/HondaAo/video-app/pkg/video/model"
	"github.com/HondaAo/video-app/utils"
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

func (r videoRepo) GetVideos(ctx context.Context, pq *utils.PaginationQuery) ([]*model.Video, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.GetVideos")
	defer span.Finish()

	rows, err := r.db.QueryContext(ctx, `SELECT video_id,title,url,category,series,end,start,level,created_at FROM videos OFFSET ? LIMIT ?`, (pq.Page-1)*pq.Size, pq.Size)
	if err != nil {
		return nil, errors.Wrap(err, "newsRepo.GetNews.QueryxContext")
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		v := &model.Video{}
		if err = rows.Scan(v.VideoId, v.Title, v.Url, v.Category, v.Series, v.End, v.Start, v.Level, v.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "GEtVideos.StructScan")
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "GetVideos.rows.Err")
	}

	return videos, nil
}

func (r videoRepo) GetVideo(ctx context.Context, id int) (*model.Video, []*model.Script, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.GetVideo")
	defer span.Finish()

	rows, err := r.db.QueryContext(ctx, `SELECT video_id,title,url,category,series,end,start,level,created_at FROM videos Where video_id = ?`, id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetVideo.QueryxContext")
	}
	defer rows.Close()

	v := &model.Video{}
	for rows.Next() {
		if err = rows.Scan(&v.VideoId, &v.Title, &v.Url, &v.Category, &v.Series, &v.End, &v.Start, &v.Level, &v.CreatedAt); err != nil {
			return nil, nil, errors.Wrap(err, "StructScan")
		}
	}

	rows, err = r.db.QueryContext(ctx, `SELECT script_id,video_id,text,ja,timestamp From scripts WHERE video_id = ?`, v.VideoId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "QueryxContext")
	}
	defer rows.Close()

	var scripts []*model.Script
	for rows.Next() {
		s := &model.Script{}
		if err = rows.Scan(&s.ScriptId, &s.VideoId, &s.Text, &s.Ja, &s.TimeStamp); err != nil {
			return nil, nil, errors.Wrap(err, "StructScan")
		}
		scripts = append(scripts, s)
	}

	return v, scripts, nil
}
