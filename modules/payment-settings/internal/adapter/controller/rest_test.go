//go:build e2e

package controller_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	paymentsettingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/controller/dto"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/testutils"
)

type PaymentSettingsControllerE2ETestSuite struct {
	suite.Suite
	pgContainer *testutils.PostgresContainer
	echo        *echo.Echo
}

func (s *PaymentSettingsControllerE2ETestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("Skipping E2E test in short mode")
	}

	s.pgContainer = testutils.SetupPostgres(s.T())
	s.pgContainer.RunMigrations(s.T(), "../../../../../migrations")

	s.echo = testutils.NewEchoForTest()
	apiGroup := s.echo.Group("/api/v1")

	paymentSettingsModule := paymentsettingsfactory.NewModule(paymentsettingsfactory.ModuleConfig{
		DB: s.pgContainer.DB,
	})

	paymentSettingsModule.RegisterHTTPHandlers(apiGroup)
}

func (s *PaymentSettingsControllerE2ETestSuite) TearDownSuite() {
	s.pgContainer.Teardown(s.T())
}

func (s *PaymentSettingsControllerE2ETestSuite) SetupTest() {
	s.pgContainer.TruncateTables(s.T(), "payment_settings", "payments")
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_CreatePaymentSetting_Success() {
	requestBody := dto.CreatePaymentSettingRequest{
		SettingKey:   "rate",
		SettingValue: "1.0",
		Currency:     "USD",
		Status:       "active",
	}

	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", requestBody)

	testutils.AssertStatusCode(s.T(), rec, http.StatusCreated)

	var response dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), rec, &response)

	assert.NotEmpty(s.T(), response.ID)
	assert.Equal(s.T(), requestBody.SettingKey, response.SettingKey)
	assert.Equal(s.T(), requestBody.SettingValue, response.SettingValue)
	assert.Equal(s.T(), requestBody.Currency, response.Currency)
	assert.Equal(s.T(), requestBody.Status, response.Status)
	assert.NotZero(s.T(), response.CreatedAt)
	assert.NotZero(s.T(), response.UpdatedAt)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_CreatePaymentSetting_InvalidJSON() {
	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", "invalid json")

	testutils.AssertStatusCode(s.T(), rec, http.StatusBadRequest)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_GetPaymentSetting_Success() {
	createReq := dto.CreatePaymentSettingRequest{
		SettingKey:   "fee",
		SettingValue: "0.5",
		Currency:     "EUR",
		Status:       "active",
	}
	createRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", createReq)
	require.Equal(s.T(), http.StatusCreated, createRec.Code)

	var createResponse dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), createRec, &createResponse)

	getRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, fmt.Sprintf("/api/v1/payment-settings/%s", createResponse.ID), nil)

	testutils.AssertStatusCode(s.T(), getRec, http.StatusOK)

	var getResponse dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), getRec, &getResponse)

	assert.Equal(s.T(), createResponse.ID, getResponse.ID)
	assert.Equal(s.T(), createResponse.SettingKey, getResponse.SettingKey)
	assert.Equal(s.T(), createResponse.SettingValue, getResponse.SettingValue)
	assert.Equal(s.T(), createResponse.Currency, getResponse.Currency)
	assert.Equal(s.T(), createResponse.Status, getResponse.Status)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_GetPaymentSetting_NotFound() {
	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payment-settings/pset_nonexistent", nil)

	testutils.AssertStatusCode(s.T(), rec, http.StatusNotFound)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_FetchPaymentSettings_Success() {
	settings := []dto.CreatePaymentSettingRequest{
		{SettingKey: "rate", SettingValue: "1.0", Currency: "USD", Status: "active"},
		{SettingKey: "fee", SettingValue: "0.5", Currency: "USD", Status: "active"},
		{SettingKey: "rate", SettingValue: "1.2", Currency: "EUR", Status: "inactive"},
	}

	for _, setting := range settings {
		rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", setting)
		require.Equal(s.T(), http.StatusCreated, rec.Code)
	}

	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payment-settings?limit=10", nil)

	testutils.AssertStatusCode(s.T(), rec, http.StatusOK)

	var response []dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), rec, &response)

	assert.Len(s.T(), response, 3)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_FetchPaymentSettings_WithFilters() {
	settings := []dto.CreatePaymentSettingRequest{
		{SettingKey: "rate", SettingValue: "1.0", Currency: "USD", Status: "active"},
		{SettingKey: "fee", SettingValue: "0.5", Currency: "USD", Status: "active"},
		{SettingKey: "rate", SettingValue: "1.2", Currency: "EUR", Status: "inactive"},
	}

	for _, setting := range settings {
		rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", setting)
		require.Equal(s.T(), http.StatusCreated, rec.Code)
	}

	tests := []struct {
		name          string
		queryParams   string
		expectedCount int
	}{
		{
			name:          "filter by currency USD",
			queryParams:   "?currency=USD&limit=10",
			expectedCount: 2,
		},
		{
			name:          "filter by settingKey rate",
			queryParams:   "?settingKey=rate&limit=10",
			expectedCount: 2,
		},
		{
			name:          "filter by status active",
			queryParams:   "?status=active&limit=10",
			expectedCount: 2,
		},
		{
			name:          "filter by multiple fields",
			queryParams:   "?currency=USD&settingKey=rate&status=active&limit=10",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			rec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payment-settings"+tt.queryParams, nil)

			testutils.AssertStatusCode(s.T(), rec, http.StatusOK)

			var response []dto.PaymentSettingResponse
			testutils.ParseJSONResponse(s.T(), rec, &response)

			assert.Len(s.T(), response, tt.expectedCount)
		})
	}
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_FetchPaymentSettings_Pagination() {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD"}
	for i := 0; i < 5; i++ {
		setting := dto.CreatePaymentSettingRequest{
			SettingKey:   "rate",
			SettingValue: "1.0",
			Currency:     currencies[i],
			Status:       "active",
		}
		rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", setting)
		require.Equal(s.T(), http.StatusCreated, rec.Code)
	}

	firstPageRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payment-settings?limit=2", nil)
	testutils.AssertStatusCode(s.T(), firstPageRec, http.StatusOK)

	var firstPage []dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), firstPageRec, &firstPage)
	assert.Len(s.T(), firstPage, 2)

	cursor := firstPageRec.Header().Get("X-Next-Cursor")
	assert.NotEmpty(s.T(), cursor)

	secondPageRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, fmt.Sprintf("/api/v1/payment-settings?limit=2&cursor=%s", cursor), nil)
	testutils.AssertStatusCode(s.T(), secondPageRec, http.StatusOK)

	var secondPage []dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), secondPageRec, &secondPage)
	assert.Len(s.T(), secondPage, 2)

	assert.NotEqual(s.T(), firstPage[0].ID, secondPage[0].ID)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_UpdatePaymentSetting_Success() {
	createReq := dto.CreatePaymentSettingRequest{
		SettingKey:   "rate",
		SettingValue: "1.0",
		Currency:     "USD",
		Status:       "active",
	}
	createRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", createReq)
	require.Equal(s.T(), http.StatusCreated, createRec.Code)

	var createResponse dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), createRec, &createResponse)

	updateReq := dto.UpdatePaymentSettingRequest{
		SettingKey:   "rate",
		SettingValue: "2.0",
		Currency:     "EUR",
		Status:       "inactive",
	}

	updateRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPut, fmt.Sprintf("/api/v1/payment-settings/%s", createResponse.ID), updateReq)

	testutils.AssertStatusCode(s.T(), updateRec, http.StatusOK)

	var updateResponse dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), updateRec, &updateResponse)

	assert.Equal(s.T(), createResponse.ID, updateResponse.ID)
	assert.Equal(s.T(), updateReq.SettingKey, updateResponse.SettingKey)
	assert.Equal(s.T(), updateReq.SettingValue, updateResponse.SettingValue)
	assert.Equal(s.T(), updateReq.Currency, updateResponse.Currency)
	assert.Equal(s.T(), updateReq.Status, updateResponse.Status)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_UpdatePaymentSetting_NotFound() {
	updateReq := dto.UpdatePaymentSettingRequest{
		SettingKey:   "rate",
		SettingValue: "2.0",
		Currency:     "EUR",
		Status:       "inactive",
	}

	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPut, "/api/v1/payment-settings/pset_nonexistent", updateReq)

	testutils.AssertStatusCode(s.T(), rec, http.StatusNotFound)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_DeletePaymentSetting_Success() {
	createReq := dto.CreatePaymentSettingRequest{
		SettingKey:   "rate",
		SettingValue: "1.0",
		Currency:     "USD",
		Status:       "active",
	}
	createRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payment-settings", createReq)
	require.Equal(s.T(), http.StatusCreated, createRec.Code)

	var createResponse dto.PaymentSettingResponse
	testutils.ParseJSONResponse(s.T(), createRec, &createResponse)

	deleteRec := testutils.MakeRequest(s.T(), s.echo, http.MethodDelete, fmt.Sprintf("/api/v1/payment-settings/%s", createResponse.ID), nil)

	testutils.AssertStatusCode(s.T(), deleteRec, http.StatusNoContent)

	getRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, fmt.Sprintf("/api/v1/payment-settings/%s", createResponse.ID), nil)
	testutils.AssertStatusCode(s.T(), getRec, http.StatusNotFound)
}

func (s *PaymentSettingsControllerE2ETestSuite) TestE2E_DeletePaymentSetting_NotFound() {
	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodDelete, "/api/v1/payment-settings/pset_nonexistent", nil)

	testutils.AssertStatusCode(s.T(), rec, http.StatusNotFound)
}

func TestE2E_PaymentSettingsControllerTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}
	suite.Run(t, new(PaymentSettingsControllerE2ETestSuite))
}
