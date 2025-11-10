package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	apperrors "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
)

func ErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	if c.Response().Committed {
		return
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		_ = c.JSON(http.StatusRequestTimeout, apperrors.ErrRequestTimeout)
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		_ = c.JSON(echoErr.Code, apperrors.EchoToHTTPError(echoErr.Code, echoErr.Message))
		return
	}

	var domainError *apperrors.Error
	if errors.As(err, &domainError) && http.StatusText(domainError.Status()) != "" {
		_ = c.JSON(domainError.Status(), domainError)
		return
	}

	log.Error().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("ip", c.RealIP()).
		Msg("Unexpected error occurred")

	_ = c.JSON(http.StatusInternalServerError, apperrors.ErrInternalServerError)
}
