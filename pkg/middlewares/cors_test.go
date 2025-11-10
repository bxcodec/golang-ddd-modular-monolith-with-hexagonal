package middlewares_test

import (
	test "net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/middlewares"
)

func TestCORS(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewRecorder()

	h := middlewares.CORS()
	e.Use(h)

	e.ServeHTTP(res, req)
	assert.Equal(t, "*", res.Header().Get("Access-Control-Allow-Origin"))
}
