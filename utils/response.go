package utils

import (
	"github.com/HondaAo/video-app/log"

	"github.com/labstack/echo/v4"
)

func LogResponseError(ctx echo.Context, logger log.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}
