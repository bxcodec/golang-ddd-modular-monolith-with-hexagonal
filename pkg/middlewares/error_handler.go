package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

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
		c.JSON(http.StatusRequestTimeout, apperrors.ErrRequestTimeout)
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		c.JSON(echoErr.Code, apperrors.EchoToHTTPError(echoErr.Code, echoErr.Message))
		return
	}

	var domainError *apperrors.Error
	if errors.As(err, &domainError) && http.StatusText(domainError.Status()) != "" {
		c.JSON(domainError.Status(), domainError)
		return
	}

	log.Printf("ERROR: Unexpected error occurred: %v", err)
	c.JSON(http.StatusInternalServerError, apperrors.ErrInternalServerError)
}
