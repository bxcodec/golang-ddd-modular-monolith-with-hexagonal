package repository

import (
	"testing"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	pkgerrors "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PaymentRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutils.PostgresContainer
	repo        *paymentRepository
}

func (s *PaymentRepositoryTestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("Skipping repository integration test in short mode")
	}
	s.pgContainer = testutils.SetupPostgres(s.T())
	s.pgContainer.RunMigrations(s.T(), "../../../../../migrations")
	s.repo = &paymentRepository{db: s.pgContainer.DB}
}

func (s *PaymentRepositoryTestSuite) TearDownSuite() {
	s.pgContainer.Teardown(s.T())
}

func (s *PaymentRepositoryTestSuite) SetupTest() {
	s.pgContainer.TruncateTables(s.T(), "payments", "payment_settings")
}

func (s *PaymentRepositoryTestSuite) TestCreatePayment() {
	tests := []struct {
		name        string
		payment     *payment.Payment
		expectError bool
	}{
		{
			name: "successful payment creation",
			payment: &payment.Payment{
				Amount:   100.50,
				Currency: "USD",
				Status:   "pending",
			},
			expectError: false,
		},
		{
			name: "create payment with different currency",
			payment: &payment.Payment{
				Amount:   250.00,
				Currency: "EUR",
				Status:   "completed",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.CreatePayment(tt.payment)

			if tt.expectError {
				assert.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)
				assert.NotEmpty(s.T(), tt.payment.ID)
				assert.NotZero(s.T(), tt.payment.CreatedAt)
				assert.NotZero(s.T(), tt.payment.UpdatedAt)
			}
		})
	}
}

func (s *PaymentRepositoryTestSuite) TestGetPayment() {
	createdPayment := &payment.Payment{
		Amount:   100.50,
		Currency: "USD",
		Status:   "pending",
	}
	err := s.repo.CreatePayment(createdPayment)
	require.NoError(s.T(), err)

	tests := []struct {
		name        string
		paymentID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "successful payment retrieval",
			paymentID:   createdPayment.ID,
			expectError: false,
		},
		{
			name:        "payment not found",
			paymentID:   "pay_nonexistent",
			expectError: true,
			expectedErr: pkgerrors.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := s.repo.GetPayment(tt.paymentID)

			if tt.expectError {
				require.Error(s.T(), err)
				if tt.expectedErr != nil {
					assert.Equal(s.T(), tt.expectedErr, err)
				}
			} else {
				require.NoError(s.T(), err)
				assert.Equal(s.T(), createdPayment.ID, result.ID)
				assert.Equal(s.T(), createdPayment.Amount, result.Amount)
				assert.Equal(s.T(), createdPayment.Currency, result.Currency)
				assert.Equal(s.T(), createdPayment.Status, result.Status)
			}
		})
	}
}

func (s *PaymentRepositoryTestSuite) TestFetchPayments() {
	payments := []*payment.Payment{
		{Amount: 100.00, Currency: "USD", Status: "pending"},
		{Amount: 200.00, Currency: "USD", Status: "completed"},
		{Amount: 300.00, Currency: "EUR", Status: "pending"},
	}

	for _, p := range payments {
		err := s.repo.CreatePayment(p)
		require.NoError(s.T(), err)
	}

	tests := []struct {
		name          string
		params        payment.FetchPaymentsParams
		expectedCount int
		expectCursor  bool
	}{
		{
			name: "fetch all payments",
			params: payment.FetchPaymentsParams{
				Limit: 10,
			},
			expectedCount: 3,
			expectCursor:  false,
		},
		{
			name: "fetch with limit",
			params: payment.FetchPaymentsParams{
				Limit: 2,
			},
			expectedCount: 2,
			expectCursor:  true,
		},
		{
			name: "filter by currency",
			params: payment.FetchPaymentsParams{
				Limit:    10,
				Currency: "USD",
			},
			expectedCount: 2,
			expectCursor:  false,
		},
		{
			name: "filter by status",
			params: payment.FetchPaymentsParams{
				Limit:  10,
				Status: "pending",
			},
			expectedCount: 2,
			expectCursor:  false,
		},
		{
			name: "filter by currency and status",
			params: payment.FetchPaymentsParams{
				Limit:    10,
				Currency: "USD",
				Status:   "completed",
			},
			expectedCount: 1,
			expectCursor:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, cursor, err := s.repo.FetchPayments(tt.params)

			require.NoError(s.T(), err)
			assert.Len(s.T(), result, tt.expectedCount)

			if tt.expectCursor {
				assert.NotEmpty(s.T(), cursor)
			} else {
				assert.Empty(s.T(), cursor)
			}
		})
	}
}

func (s *PaymentRepositoryTestSuite) TestFetchPayments_Pagination() {
	for i := 0; i < 5; i++ {
		p := &payment.Payment{
			Amount:   float64(100 * (i + 1)),
			Currency: "USD",
			Status:   "pending",
		}
		err := s.repo.CreatePayment(p)
		require.NoError(s.T(), err)
	}

	firstPage, cursor, err := s.repo.FetchPayments(payment.FetchPaymentsParams{Limit: 2})
	require.NoError(s.T(), err)
	assert.Len(s.T(), firstPage, 2)
	assert.NotEmpty(s.T(), cursor)

	secondPage, cursor2, err := s.repo.FetchPayments(payment.FetchPaymentsParams{
		Limit:  2,
		Cursor: cursor,
	})
	require.NoError(s.T(), err)
	assert.Len(s.T(), secondPage, 2)
	assert.NotEmpty(s.T(), cursor2)

	assert.NotEqual(s.T(), firstPage[0].ID, secondPage[0].ID)
	assert.NotEqual(s.T(), firstPage[1].ID, secondPage[1].ID)
}

func (s *PaymentRepositoryTestSuite) TestUpdatePayment() {
	createdPayment := &payment.Payment{
		Amount:   100.50,
		Currency: "USD",
		Status:   "pending",
	}
	err := s.repo.CreatePayment(createdPayment)
	require.NoError(s.T(), err)

	tests := []struct {
		name        string
		payment     *payment.Payment
		expectError bool
		expectedErr error
	}{
		{
			name: "successful payment update",
			payment: &payment.Payment{
				ID:       createdPayment.ID,
				Amount:   150.00,
				Currency: "EUR",
				Status:   "completed",
			},
			expectError: false,
		},
		{
			name: "payment not found",
			payment: &payment.Payment{
				ID:       "pay_nonexistent",
				Amount:   150.00,
				Currency: "EUR",
				Status:   "completed",
			},
			expectError: true,
			expectedErr: pkgerrors.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.UpdatePayment(tt.payment)

			if tt.expectError {
				require.Error(s.T(), err)
				if tt.expectedErr != nil {
					assert.Equal(s.T(), tt.expectedErr, err)
				}
			} else {
				require.NoError(s.T(), err)

				updated, err := s.repo.GetPayment(tt.payment.ID)
				require.NoError(s.T(), err)
				assert.Equal(s.T(), tt.payment.Amount, updated.Amount)
				assert.Equal(s.T(), tt.payment.Currency, updated.Currency)
				assert.Equal(s.T(), tt.payment.Status, updated.Status)
			}
		})
	}
}

func (s *PaymentRepositoryTestSuite) TestDeletePayment() {
	createdPayment := &payment.Payment{
		Amount:   100.50,
		Currency: "USD",
		Status:   "pending",
	}
	err := s.repo.CreatePayment(createdPayment)
	require.NoError(s.T(), err)

	tests := []struct {
		name        string
		paymentID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "successful payment deletion",
			paymentID:   createdPayment.ID,
			expectError: false,
		},
		{
			name:        "payment not found",
			paymentID:   "pay_nonexistent",
			expectError: true,
			expectedErr: pkgerrors.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.DeletePayment(tt.paymentID)

			if tt.expectError {
				require.Error(s.T(), err)
				if tt.expectedErr != nil {
					assert.Equal(s.T(), tt.expectedErr, err)
				}
			} else {
				require.NoError(s.T(), err)

				_, err := s.repo.GetPayment(tt.paymentID)
				assert.Equal(s.T(), pkgerrors.ErrDataNotFound, err)
			}
		})
	}
}

func TestPaymentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}
