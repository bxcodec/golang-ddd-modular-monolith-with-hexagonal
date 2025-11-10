package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports/mocks"
	pkgerrors "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
)

func TestPaymentService_CreatePayment(t *testing.T) {
	tests := []struct {
		name                  string
		payment               *payment.Payment
		mockSettingsResponse  []paymentsettings.PaymentSetting
		mockSettingsCursor    string
		mockSettingsError     error
		mockCreateError       error
		expectError           bool
		expectedErrorContains string
	}{
		{
			name: "successful payment creation",
			payment: &payment.Payment{
				Amount:   100.50,
				Currency: "USD",
				Status:   "pending",
			},
			mockSettingsResponse: []paymentsettings.PaymentSetting{
				{
					ID:           "pset_123",
					SettingKey:   "rate",
					SettingValue: "1.0",
					Currency:     "USD",
					Status:       "active",
				},
			},
			mockSettingsCursor: "",
			mockSettingsError:  nil,
			mockCreateError:    nil,
			expectError:        false,
		},
		{
			name: "payment settings fetch error",
			payment: &payment.Payment{
				Amount:   100.50,
				Currency: "EUR",
				Status:   "pending",
			},
			mockSettingsResponse:  nil,
			mockSettingsCursor:    "",
			mockSettingsError:     errors.New("database connection failed"),
			mockCreateError:       nil,
			expectError:           true,
			expectedErrorContains: "database connection failed",
		},
		{
			name: "payment creation error",
			payment: &payment.Payment{
				Amount:   100.50,
				Currency: "USD",
				Status:   "pending",
			},
			mockSettingsResponse: []paymentsettings.PaymentSetting{
				{
					ID:       "pset_123",
					Currency: "USD",
				},
			},
			mockSettingsCursor:    "",
			mockSettingsError:     nil,
			mockCreateError:       pkgerrors.ErrDuplicatedData,
			expectError:           true,
			expectedErrorContains: "DATA_DUPLICATE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentRepository(t)
			mockSettingsPort := mocks.NewMockIPaymentSettingsPort(t)

			mockSettingsPort.On("FetchPaymentSettings", mock.MatchedBy(func(params paymentsettings.PaymentSettingFetchParams) bool {
				return params.Currency == tt.payment.Currency && params.Limit == 1
			})).Return(tt.mockSettingsResponse, tt.mockSettingsCursor, tt.mockSettingsError)

			if tt.mockSettingsError == nil {
				mockRepo.On("CreatePayment", tt.payment).Return(tt.mockCreateError)
			}

			service := NewPaymentService(mockRepo, mockSettingsPort)
			err := service.CreatePayment(tt.payment)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorContains != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaymentService_GetPayment(t *testing.T) {
	tests := []struct {
		name          string
		paymentID     string
		mockPayment   payment.Payment
		mockError     error
		expectError   bool
		expectedError error
	}{
		{
			name:      "successful payment retrieval",
			paymentID: "pay_123",
			mockPayment: payment.Payment{
				ID:       "pay_123",
				Amount:   100.50,
				Currency: "USD",
				Status:   "completed",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:          "payment not found",
			paymentID:     "pay_nonexistent",
			mockPayment:   payment.Payment{},
			mockError:     pkgerrors.ErrDataNotFound,
			expectError:   true,
			expectedError: pkgerrors.ErrDataNotFound,
		},
		{
			name:        "database error",
			paymentID:   "pay_123",
			mockPayment: payment.Payment{},
			mockError:   errors.New("connection timeout"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentRepository(t)
			mockSettingsPort := mocks.NewMockIPaymentSettingsPort(t)

			mockRepo.On("GetPayment", tt.paymentID).Return(tt.mockPayment, tt.mockError)

			service := NewPaymentService(mockRepo, mockSettingsPort)
			result, err := service.GetPayment(tt.paymentID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPayment, result)
			}
		})
	}
}

func TestPaymentService_FetchPayments(t *testing.T) {
	tests := []struct {
		name         string
		params       payment.FetchPaymentsParams
		mockPayments []payment.Payment
		mockCursor   string
		mockError    error
		expectError  bool
	}{
		{
			name: "successful fetch with results",
			params: payment.FetchPaymentsParams{
				Cursor:   "",
				Limit:    10,
				Currency: "USD",
				Status:   "completed",
			},
			mockPayments: []payment.Payment{
				{
					ID:        "pay_123",
					Amount:    100.50,
					Currency:  "USD",
					Status:    "completed",
					CreatedAt: time.Now(),
				},
				{
					ID:        "pay_124",
					Amount:    200.00,
					Currency:  "USD",
					Status:    "completed",
					CreatedAt: time.Now(),
				},
			},
			mockCursor:  "encoded_cursor",
			mockError:   nil,
			expectError: false,
		},
		{
			name: "successful fetch with no results",
			params: payment.FetchPaymentsParams{
				Cursor: "",
				Limit:  10,
			},
			mockPayments: []payment.Payment{},
			mockCursor:   "",
			mockError:    nil,
			expectError:  false,
		},
		{
			name: "database error",
			params: payment.FetchPaymentsParams{
				Cursor: "",
				Limit:  10,
			},
			mockPayments: nil,
			mockCursor:   "",
			mockError:    errors.New("query failed"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentRepository(t)
			mockSettingsPort := mocks.NewMockIPaymentSettingsPort(t)

			mockRepo.On("FetchPayments", tt.params).Return(tt.mockPayments, tt.mockCursor, tt.mockError)

			service := NewPaymentService(mockRepo, mockSettingsPort)
			result, cursor, err := service.FetchPayments(tt.params)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPayments, result)
				assert.Equal(t, tt.mockCursor, cursor)
			}
		})
	}
}

func TestPaymentService_UpdatePayment(t *testing.T) {
	tests := []struct {
		name        string
		payment     *payment.Payment
		mockError   error
		expectError bool
	}{
		{
			name: "successful payment update",
			payment: &payment.Payment{
				ID:       "pay_123",
				Amount:   150.00,
				Currency: "USD",
				Status:   "completed",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "payment not found",
			payment: &payment.Payment{
				ID:       "pay_nonexistent",
				Amount:   150.00,
				Currency: "USD",
				Status:   "completed",
			},
			mockError:   pkgerrors.ErrDataNotFound,
			expectError: true,
		},
		{
			name: "database error",
			payment: &payment.Payment{
				ID:       "pay_123",
				Amount:   150.00,
				Currency: "USD",
				Status:   "completed",
			},
			mockError:   errors.New("update failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentRepository(t)
			mockSettingsPort := mocks.NewMockIPaymentSettingsPort(t)

			mockRepo.On("UpdatePayment", tt.payment).Return(tt.mockError)

			service := NewPaymentService(mockRepo, mockSettingsPort)
			err := service.UpdatePayment(tt.payment)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaymentService_DeletePayment(t *testing.T) {
	tests := []struct {
		name        string
		paymentID   string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful payment deletion",
			paymentID:   "pay_123",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "payment not found",
			paymentID:   "pay_nonexistent",
			mockError:   pkgerrors.ErrDataNotFound,
			expectError: true,
		},
		{
			name:        "database error",
			paymentID:   "pay_123",
			mockError:   errors.New("delete failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentRepository(t)
			mockSettingsPort := mocks.NewMockIPaymentSettingsPort(t)

			mockRepo.On("DeletePayment", tt.paymentID).Return(tt.mockError)

			service := NewPaymentService(mockRepo, mockSettingsPort)
			err := service.DeletePayment(tt.paymentID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
