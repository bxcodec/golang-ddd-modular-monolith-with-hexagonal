package service

import (
	"errors"
	"testing"
	"time"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/ports/mocks"
	pkgerrors "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentSettingsService_CreatePaymentSetting(t *testing.T) {
	tests := []struct {
		name                  string
		setting               *paymentsettings.PaymentSetting
		mockError             error
		expectError           bool
		expectedErrorContains string
	}{
		{
			name: "successful payment setting creation",
			setting: &paymentsettings.PaymentSetting{
				SettingKey:   "rate",
				SettingValue: "1.0",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "duplicate payment setting",
			setting: &paymentsettings.PaymentSetting{
				SettingKey:   "rate",
				SettingValue: "1.0",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:             pkgerrors.ErrDuplicatedData,
			expectError:           true,
			expectedErrorContains: "DATA_DUPLICATE",
		},
		{
			name: "database error",
			setting: &paymentsettings.PaymentSetting{
				SettingKey:   "rate",
				SettingValue: "1.0",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:   errors.New("connection timeout"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentSettingsRepository(t)

			mockRepo.On("CreatePaymentSetting", tt.setting).Return(tt.mockError)

			service := NewPaymentSettingsService(mockRepo)
			err := service.CreatePaymentSetting(tt.setting)

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

func TestPaymentSettingsService_GetPaymentSetting(t *testing.T) {
	tests := []struct {
		name          string
		settingID     string
		mockSetting   paymentsettings.PaymentSetting
		mockError     error
		expectError   bool
		expectedError error
	}{
		{
			name:      "successful payment setting retrieval",
			settingID: "pset_123",
			mockSetting: paymentsettings.PaymentSetting{
				ID:           "pset_123",
				SettingKey:   "rate",
				SettingValue: "1.0",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:          "payment setting not found",
			settingID:     "pset_nonexistent",
			mockSetting:   paymentsettings.PaymentSetting{},
			mockError:     pkgerrors.ErrDataNotFound,
			expectError:   true,
			expectedError: pkgerrors.ErrDataNotFound,
		},
		{
			name:        "database error",
			settingID:   "pset_123",
			mockSetting: paymentsettings.PaymentSetting{},
			mockError:   errors.New("connection timeout"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentSettingsRepository(t)

			mockRepo.On("GetPaymentSetting", tt.settingID).Return(tt.mockSetting, tt.mockError)

			service := NewPaymentSettingsService(mockRepo)
			result, err := service.GetPaymentSetting(tt.settingID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockSetting, result)
			}
		})
	}
}

func TestPaymentSettingsService_FetchPaymentSettings(t *testing.T) {
	tests := []struct {
		name         string
		params       paymentsettings.PaymentSettingFetchParams
		mockSettings []paymentsettings.PaymentSetting
		mockCursor   string
		mockError    error
		expectError  bool
	}{
		{
			name: "successful fetch with results",
			params: paymentsettings.PaymentSettingFetchParams{
				Cursor:     "",
				Limit:      10,
				Currency:   "USD",
				SettingKey: "rate",
				Status:     "active",
			},
			mockSettings: []paymentsettings.PaymentSetting{
				{
					ID:           "pset_123",
					SettingKey:   "rate",
					SettingValue: "1.0",
					Currency:     "USD",
					Status:       "active",
					CreatedAt:    time.Now(),
				},
				{
					ID:           "pset_124",
					SettingKey:   "rate",
					SettingValue: "1.1",
					Currency:     "USD",
					Status:       "active",
					CreatedAt:    time.Now(),
				},
			},
			mockCursor:  "encoded_cursor",
			mockError:   nil,
			expectError: false,
		},
		{
			name: "successful fetch with no results",
			params: paymentsettings.PaymentSettingFetchParams{
				Cursor: "",
				Limit:  10,
			},
			mockSettings: []paymentsettings.PaymentSetting{},
			mockCursor:   "",
			mockError:    nil,
			expectError:  false,
		},
		{
			name: "database error",
			params: paymentsettings.PaymentSettingFetchParams{
				Cursor: "",
				Limit:  10,
			},
			mockSettings: nil,
			mockCursor:   "",
			mockError:    errors.New("query failed"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentSettingsRepository(t)

			mockRepo.On("FetchPaymentSettings", mock.MatchedBy(func(params paymentsettings.PaymentSettingFetchParams) bool {
				return params.Cursor == tt.params.Cursor &&
					params.Limit == tt.params.Limit &&
					params.Currency == tt.params.Currency &&
					params.SettingKey == tt.params.SettingKey &&
					params.Status == tt.params.Status
			})).Return(tt.mockSettings, tt.mockCursor, tt.mockError)

			service := NewPaymentSettingsService(mockRepo)
			result, cursor, err := service.FetchPaymentSettings(tt.params)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockSettings, result)
				assert.Equal(t, tt.mockCursor, cursor)
			}
		})
	}
}

func TestPaymentSettingsService_UpdatePaymentSetting(t *testing.T) {
	tests := []struct {
		name        string
		setting     *paymentsettings.PaymentSetting
		mockError   error
		expectError bool
	}{
		{
			name: "successful payment setting update",
			setting: &paymentsettings.PaymentSetting{
				ID:           "pset_123",
				SettingKey:   "rate",
				SettingValue: "1.5",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "payment setting not found",
			setting: &paymentsettings.PaymentSetting{
				ID:           "pset_nonexistent",
				SettingKey:   "rate",
				SettingValue: "1.5",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:   pkgerrors.ErrDataNotFound,
			expectError: true,
		},
		{
			name: "database error",
			setting: &paymentsettings.PaymentSetting{
				ID:           "pset_123",
				SettingKey:   "rate",
				SettingValue: "1.5",
				Currency:     "USD",
				Status:       "active",
			},
			mockError:   errors.New("update failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentSettingsRepository(t)

			mockRepo.On("UpdatePaymentSetting", tt.setting).Return(tt.mockError)

			service := NewPaymentSettingsService(mockRepo)
			err := service.UpdatePaymentSetting(tt.setting)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaymentSettingsService_DeletePaymentSetting(t *testing.T) {
	tests := []struct {
		name        string
		settingID   string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful payment setting deletion",
			settingID:   "pset_123",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "payment setting not found",
			settingID:   "pset_nonexistent",
			mockError:   pkgerrors.ErrDataNotFound,
			expectError: true,
		},
		{
			name:        "database error",
			settingID:   "pset_123",
			mockError:   errors.New("delete failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIPaymentSettingsRepository(t)

			mockRepo.On("DeletePaymentSetting", tt.settingID).Return(tt.mockError)

			service := NewPaymentSettingsService(mockRepo)
			err := service.DeletePaymentSetting(tt.settingID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
