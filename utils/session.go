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
	SessionID string    `json:"session_id" redis:"session_id"`
	UserID    uuid.UUID `json:"user_id" redis:"user_id"`
}

// Create session in redis
func CreateSession(ctx context.Context, sess Session, expire int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRepo.CreateSession")
	defer span.Finish()

	sess.SessionID = uuid.New().String()
	sessionKey := createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", errors.WithMessage(err, "sessionRepo.CreateSession.json.Marshal")
	}
	if err = redis.Client.Set(redis.Client{}, ctx, sessionKey, string(sessBytes), time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}
	return sessionKey, nil
}

func createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, sessionID)
}
