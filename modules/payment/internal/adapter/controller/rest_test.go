//go:build e2e

package controller_test

import (
	"fmt"
	"net/http"
	"testing"

	paymentsettingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/controller/dto"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/testutils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PaymentControllerE2ETestSuite struct {
	suite.Suite
	pgContainer *testutils.PostgresContainer
	echo        *echo.Echo
}

func (s *PaymentControllerE2ETestSuite) SetupSuite() {
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

	paymentModule := factory.NewModule(factory.ModuleConfig{
		DB:                  s.pgContainer.DB,
		PaymentSettingsPort: paymentSettingsModule.Service,
	})

	paymentSettingsModule.RegisterHTTPHandlers(apiGroup)
	paymentModule.RegisterHTTPHandlers(apiGroup)
}

func (s *PaymentControllerE2ETestSuite) TearDownSuite() {
	s.pgContainer.Teardown(s.T())
}

func (s *PaymentControllerE2ETestSuite) SetupTest() {
	s.pgContainer.TruncateTables(s.T(), "payments", "payment_settings")
}

func (s *PaymentControllerE2ETestSuite) TestE2E_CreatePayment_Success() {
	requestBody := dto.CreatePaymentRequest{
		Amount:   100.50,
		Currency: "USD",
		Status:   "pending",
	}

	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", requestBody)

	testutils.AssertStatusCode(s.T(), rec, http.StatusCreated)

	var response dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), rec, &response)

	assert.NotEmpty(s.T(), response.ID)
	assert.Equal(s.T(), requestBody.Amount, response.Amount)
	assert.Equal(s.T(), requestBody.Currency, response.Currency)
	assert.Equal(s.T(), requestBody.Status, response.Status)
	assert.NotZero(s.T(), response.CreatedAt)
	assert.NotZero(s.T(), response.UpdatedAt)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_CreatePayment_InvalidJSON() {
	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", "invalid json")

	testutils.AssertStatusCode(s.T(), rec, http.StatusBadRequest)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_GetPayment_Success() {
	createReq := dto.CreatePaymentRequest{
		Amount:   150.00,
		Currency: "EUR",
		Status:   "completed",
	}
	createRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", createReq)
	require.Equal(s.T(), http.StatusCreated, createRec.Code)

	var createResponse dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), createRec, &createResponse)

	getRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, fmt.Sprintf("/api/v1/payments/%s", createResponse.ID), nil)

	testutils.AssertStatusCode(s.T(), getRec, http.StatusOK)

	var getResponse dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), getRec, &getResponse)

	assert.Equal(s.T(), createResponse.ID, getResponse.ID)
	assert.Equal(s.T(), createResponse.Amount, getResponse.Amount)
	assert.Equal(s.T(), createResponse.Currency, getResponse.Currency)
	assert.Equal(s.T(), createResponse.Status, getResponse.Status)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_GetPayment_NotFound() {
	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payments/pay_nonexistent", nil)

	testutils.AssertStatusCode(s.T(), rec, http.StatusNotFound)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_FetchPayments_Success() {
	payments := []dto.CreatePaymentRequest{
		{Amount: 100.00, Currency: "USD", Status: "pending"},
		{Amount: 200.00, Currency: "USD", Status: "completed"},
		{Amount: 300.00, Currency: "EUR", Status: "pending"},
	}

	for _, p := range payments {
		rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", p)
		require.Equal(s.T(), http.StatusCreated, rec.Code)
	}

	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payments?limit=10", nil)

	testutils.AssertStatusCode(s.T(), rec, http.StatusOK)

	var response []dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), rec, &response)

	assert.Len(s.T(), response, 3)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_FetchPayments_WithFilters() {
	payments := []dto.CreatePaymentRequest{
		{Amount: 100.00, Currency: "USD", Status: "pending"},
		{Amount: 200.00, Currency: "USD", Status: "completed"},
		{Amount: 300.00, Currency: "EUR", Status: "pending"},
	}

	for _, p := range payments {
		rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", p)
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
			name:          "filter by status pending",
			queryParams:   "?status=pending&limit=10",
			expectedCount: 2,
		},
		{
			name:          "filter by currency and status",
			queryParams:   "?currency=USD&status=completed&limit=10",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			rec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payments"+tt.queryParams, nil)

			testutils.AssertStatusCode(s.T(), rec, http.StatusOK)

			var response []dto.PaymentResponse
			testutils.ParseJSONResponse(s.T(), rec, &response)

			assert.Len(s.T(), response, tt.expectedCount)
		})
	}
}

func (s *PaymentControllerE2ETestSuite) TestE2E_FetchPayments_Pagination() {
	for i := 0; i < 5; i++ {
		p := dto.CreatePaymentRequest{
			Amount:   float64(100 * (i + 1)),
			Currency: "USD",
			Status:   "pending",
		}
		rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", p)
		require.Equal(s.T(), http.StatusCreated, rec.Code)
	}

	firstPageRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, "/api/v1/payments?limit=2", nil)
	testutils.AssertStatusCode(s.T(), firstPageRec, http.StatusOK)

	var firstPage []dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), firstPageRec, &firstPage)
	assert.Len(s.T(), firstPage, 2)

	cursor := firstPageRec.Header().Get("X-Next-Cursor")
	assert.NotEmpty(s.T(), cursor)

	secondPageRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, fmt.Sprintf("/api/v1/payments?limit=2&cursor=%s", cursor), nil)
	testutils.AssertStatusCode(s.T(), secondPageRec, http.StatusOK)

	var secondPage []dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), secondPageRec, &secondPage)
	assert.Len(s.T(), secondPage, 2)

	assert.NotEqual(s.T(), firstPage[0].ID, secondPage[0].ID)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_UpdatePayment_Success() {
	createReq := dto.CreatePaymentRequest{
		Amount:   100.00,
		Currency: "USD",
		Status:   "pending",
	}
	createRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", createReq)
	require.Equal(s.T(), http.StatusCreated, createRec.Code)

	var createResponse dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), createRec, &createResponse)

	updateReq := dto.UpdatePaymentRequest{
		Amount:   200.00,
		Currency: "EUR",
		Status:   "completed",
	}

	updateRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPut, fmt.Sprintf("/api/v1/payments/%s", createResponse.ID), updateReq)

	testutils.AssertStatusCode(s.T(), updateRec, http.StatusOK)

	var updateResponse dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), updateRec, &updateResponse)

	assert.Equal(s.T(), createResponse.ID, updateResponse.ID)
	assert.Equal(s.T(), updateReq.Amount, updateResponse.Amount)
	assert.Equal(s.T(), updateReq.Currency, updateResponse.Currency)
	assert.Equal(s.T(), updateReq.Status, updateResponse.Status)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_UpdatePayment_NotFound() {
	updateReq := dto.UpdatePaymentRequest{
		Amount:   200.00,
		Currency: "EUR",
		Status:   "completed",
	}

	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodPut, "/api/v1/payments/pay_nonexistent", updateReq)

	testutils.AssertStatusCode(s.T(), rec, http.StatusNotFound)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_DeletePayment_Success() {
	createReq := dto.CreatePaymentRequest{
		Amount:   100.00,
		Currency: "USD",
		Status:   "pending",
	}
	createRec := testutils.MakeRequest(s.T(), s.echo, http.MethodPost, "/api/v1/payments", createReq)
	require.Equal(s.T(), http.StatusCreated, createRec.Code)

	var createResponse dto.PaymentResponse
	testutils.ParseJSONResponse(s.T(), createRec, &createResponse)

	deleteRec := testutils.MakeRequest(s.T(), s.echo, http.MethodDelete, fmt.Sprintf("/api/v1/payments/%s", createResponse.ID), nil)

	testutils.AssertStatusCode(s.T(), deleteRec, http.StatusNoContent)

	getRec := testutils.MakeRequest(s.T(), s.echo, http.MethodGet, fmt.Sprintf("/api/v1/payments/%s", createResponse.ID), nil)
	testutils.AssertStatusCode(s.T(), getRec, http.StatusNotFound)
}

func (s *PaymentControllerE2ETestSuite) TestE2E_DeletePayment_NotFound() {
	rec := testutils.MakeRequest(s.T(), s.echo, http.MethodDelete, "/api/v1/payments/pay_nonexistent", nil)

	testutils.AssertStatusCode(s.T(), rec, http.StatusNotFound)
}

func TestE2E_PaymentControllerTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}
	suite.Run(t, new(PaymentControllerE2ETestSuite))
}
