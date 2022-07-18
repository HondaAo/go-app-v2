package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	basePrefix = "api-session:"
)

type Session struct {
	SessionID string `json:"session_id" redis:"session_id"`
	UserID    string `json:"user_id" redis:"user_id"`
}

type UserCtxKey struct{}

// Create session in redis
func CreateSession(ctx context.Context, redisClient redis.Client, sess Session, expire int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRepo.CreateSession")
	defer span.Finish()

	sess.SessionID = uuid.New().String()
	sessionKey := createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", errors.WithMessage(err, "sessionRepo.CreateSession.json.Marshal")
	}
	if err = redisClient.Set(ctx, sessionKey, string(sessBytes), time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}
	return sessionKey, nil
}

// get session by id
func GetSessionByID(ctx context.Context, redis redis.Client, sessionID string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.GetSessionByID")
	defer span.Finish()

	sessBytes, err := redis.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRep.GetSessionByID.redisClient.Get")
	}

	sess := &Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.json.Unmarshal")
	}
	return sess, nil
}

func createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, sessionID)
}
