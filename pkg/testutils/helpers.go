package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func MakeRequest(t *testing.T, e *echo.Echo, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}

func ParseJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, v interface{}) {
	err := json.Unmarshal(rec.Body.Bytes(), v)
	require.NoError(t, err)
}

func AssertStatusCode(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int) {
	require.Equal(t, expectedStatus, rec.Code, "Response body: %s", rec.Body.String())
}

func AssertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	AssertStatusCode(t, rec, expectedStatus)
	if expectedBody != nil {
		actualJSON := make(map[string]interface{})
		ParseJSONResponse(t, rec, &actualJSON)

		expectedJSON, err := json.Marshal(expectedBody)
		require.NoError(t, err)

		expectedMap := make(map[string]interface{})
		err = json.Unmarshal(expectedJSON, &expectedMap)
		require.NoError(t, err)

		require.Equal(t, expectedMap, actualJSON)
	}
}

func NewEchoForTest() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		} else {
			type statusCoder interface {
				Status() int
			}
			if sc, ok := err.(statusCoder); ok {
				code = sc.Status()
			}
		}

		c.JSON(code, map[string]string{"error": message})
	}
	return e
}
