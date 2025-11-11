package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	pkgerrors "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/testutils"
)

type PaymentSettingsRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutils.PostgresContainer
	repo        *PaymentSettingsRepository
}

func (s *PaymentSettingsRepositoryTestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("Skipping repository integration test in short mode")
	}
	s.pgContainer = testutils.SetupPostgres(s.T())
	s.pgContainer.RunMigrations(s.T(), "../../../../../migrations")
	s.repo = &PaymentSettingsRepository{db: s.pgContainer.DB}
}

func (s *PaymentSettingsRepositoryTestSuite) TearDownSuite() {
	s.pgContainer.Teardown(s.T())
}

func (s *PaymentSettingsRepositoryTestSuite) SetupTest() {
	s.pgContainer.TruncateTables(s.T(), "payment_settings_module.payment_settings", "payment_module.payments")
}

func (s *PaymentSettingsRepositoryTestSuite) TestCreatePaymentSetting() {
	tests := []struct {
		name        string
		setting     *paymentsettings.PaymentSetting
		expectError bool
	}{
		{
			name: "successful payment setting creation",
			setting: &paymentsettings.PaymentSetting{
				SettingKey:   "rate",
				SettingValue: "1.0",
				Currency:     "USD",
				Status:       "active",
			},
			expectError: false,
		},
		{
			name: "create payment setting with different currency",
			setting: &paymentsettings.PaymentSetting{
				SettingKey:   "fee",
				SettingValue: "0.5",
				Currency:     "EUR",
				Status:       "active",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.CreatePaymentSetting(tt.setting)

			if tt.expectError {
				assert.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)
				assert.NotEmpty(s.T(), tt.setting.ID)
				assert.NotZero(s.T(), tt.setting.CreatedAt)
				assert.NotZero(s.T(), tt.setting.UpdatedAt)
			}
		})
	}
}

func (s *PaymentSettingsRepositoryTestSuite) TestGetPaymentSetting() {
	createdSetting := &paymentsettings.PaymentSetting{
		SettingKey:   "rate",
		SettingValue: "1.0",
		Currency:     "USD",
		Status:       "active",
	}
	err := s.repo.CreatePaymentSetting(createdSetting)
	require.NoError(s.T(), err)

	tests := []struct {
		name        string
		settingID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "successful payment setting retrieval",
			settingID:   createdSetting.ID,
			expectError: false,
		},
		{
			name:        "payment setting not found",
			settingID:   "pset_nonexistent",
			expectError: true,
			expectedErr: pkgerrors.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := s.repo.GetPaymentSetting(tt.settingID)

			if tt.expectError {
				require.Error(s.T(), err)
				if tt.expectedErr != nil {
					assert.Equal(s.T(), tt.expectedErr, err)
				}
			} else {
				require.NoError(s.T(), err)
				assert.Equal(s.T(), createdSetting.ID, result.ID)
				assert.Equal(s.T(), createdSetting.SettingKey, result.SettingKey)
				assert.Equal(s.T(), createdSetting.SettingValue, result.SettingValue)
				assert.Equal(s.T(), createdSetting.Currency, result.Currency)
				assert.Equal(s.T(), createdSetting.Status, result.Status)
			}
		})
	}
}

func (s *PaymentSettingsRepositoryTestSuite) TestFetchPaymentSettings() {
	settings := []*paymentsettings.PaymentSetting{
		{SettingKey: "rate", SettingValue: "1.0", Currency: "USD", Status: "active"},
		{SettingKey: "fee", SettingValue: "0.5", Currency: "USD", Status: "active"},
		{SettingKey: "rate", SettingValue: "1.2", Currency: "EUR", Status: "inactive"},
	}

	for _, setting := range settings {
		err := s.repo.CreatePaymentSetting(setting)
		require.NoError(s.T(), err)
	}

	tests := []struct {
		name          string
		params        paymentsettings.PaymentSettingFetchParams
		expectedCount int
		expectCursor  bool
	}{
		{
			name: "fetch all payment settings",
			params: paymentsettings.PaymentSettingFetchParams{
				Limit: 10,
			},
			expectedCount: 3,
			expectCursor:  false,
		},
		{
			name: "fetch with limit",
			params: paymentsettings.PaymentSettingFetchParams{
				Limit: 2,
			},
			expectedCount: 2,
			expectCursor:  true,
		},
		{
			name: "filter by currency",
			params: paymentsettings.PaymentSettingFetchParams{
				Limit:    10,
				Currency: "USD",
			},
			expectedCount: 2,
			expectCursor:  false,
		},
		{
			name: "filter by setting key",
			params: paymentsettings.PaymentSettingFetchParams{
				Limit:      10,
				SettingKey: "rate",
			},
			expectedCount: 2,
			expectCursor:  false,
		},
		{
			name: "filter by status",
			params: paymentsettings.PaymentSettingFetchParams{
				Limit:  10,
				Status: "active",
			},
			expectedCount: 2,
			expectCursor:  false,
		},
		{
			name: "filter by multiple fields",
			params: paymentsettings.PaymentSettingFetchParams{
				Limit:      10,
				Currency:   "USD",
				SettingKey: "rate",
				Status:     "active",
			},
			expectedCount: 1,
			expectCursor:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, cursor, err := s.repo.FetchPaymentSettings(tt.params)

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

func (s *PaymentSettingsRepositoryTestSuite) TestFetchPaymentSettings_Pagination() {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD"}
	for i := 0; i < 5; i++ {
		setting := &paymentsettings.PaymentSetting{
			SettingKey:   "rate",
			SettingValue: "1.0",
			Currency:     currencies[i],
			Status:       "active",
		}
		err := s.repo.CreatePaymentSetting(setting)
		require.NoError(s.T(), err)
	}

	firstPage, cursor, err := s.repo.FetchPaymentSettings(paymentsettings.PaymentSettingFetchParams{Limit: 2})
	require.NoError(s.T(), err)
	assert.Len(s.T(), firstPage, 2)
	assert.NotEmpty(s.T(), cursor)

	secondPage, cursor2, err := s.repo.FetchPaymentSettings(paymentsettings.PaymentSettingFetchParams{
		Limit:  2,
		Cursor: cursor,
	})
	require.NoError(s.T(), err)
	assert.Len(s.T(), secondPage, 2)
	assert.NotEmpty(s.T(), cursor2)

	assert.NotEqual(s.T(), firstPage[0].ID, secondPage[0].ID)
	assert.NotEqual(s.T(), firstPage[1].ID, secondPage[1].ID)
}

func (s *PaymentSettingsRepositoryTestSuite) TestUpdatePaymentSetting() {
	createdSetting := &paymentsettings.PaymentSetting{
		SettingKey:   "rate",
		SettingValue: "1.0",
		Currency:     "USD",
		Status:       "active",
	}
	err := s.repo.CreatePaymentSetting(createdSetting)
	require.NoError(s.T(), err)

	tests := []struct {
		name        string
		setting     *paymentsettings.PaymentSetting
		expectError bool
		expectedErr error
	}{
		{
			name: "successful payment setting update",
			setting: &paymentsettings.PaymentSetting{
				ID:           createdSetting.ID,
				SettingKey:   "rate",
				SettingValue: "1.5",
				Currency:     "EUR",
				Status:       "inactive",
			},
			expectError: false,
		},
		{
			name: "payment setting not found",
			setting: &paymentsettings.PaymentSetting{
				ID:           "pset_nonexistent",
				SettingKey:   "rate",
				SettingValue: "1.5",
				Currency:     "EUR",
				Status:       "inactive",
			},
			expectError: true,
			expectedErr: pkgerrors.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.UpdatePaymentSetting(tt.setting)

			if tt.expectError {
				require.Error(s.T(), err)
				if tt.expectedErr != nil {
					assert.Equal(s.T(), tt.expectedErr, err)
				}
			} else {
				require.NoError(s.T(), err)

				updated, err := s.repo.GetPaymentSetting(tt.setting.ID)
				require.NoError(s.T(), err)
				assert.Equal(s.T(), tt.setting.SettingKey, updated.SettingKey)
				assert.Equal(s.T(), tt.setting.SettingValue, updated.SettingValue)
				assert.Equal(s.T(), tt.setting.Currency, updated.Currency)
				assert.Equal(s.T(), tt.setting.Status, updated.Status)
			}
		})
	}
}

func (s *PaymentSettingsRepositoryTestSuite) TestDeletePaymentSetting() {
	createdSetting := &paymentsettings.PaymentSetting{
		SettingKey:   "rate",
		SettingValue: "1.0",
		Currency:     "USD",
		Status:       "active",
	}
	err := s.repo.CreatePaymentSetting(createdSetting)
	require.NoError(s.T(), err)

	tests := []struct {
		name        string
		settingID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "successful payment setting deletion",
			settingID:   createdSetting.ID,
			expectError: false,
		},
		{
			name:        "payment setting not found",
			settingID:   "pset_nonexistent",
			expectError: true,
			expectedErr: pkgerrors.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.DeletePaymentSetting(tt.settingID)

			if tt.expectError {
				require.Error(s.T(), err)
				if tt.expectedErr != nil {
					assert.Equal(s.T(), tt.expectedErr, err)
				}
			} else {
				require.NoError(s.T(), err)

				_, err := s.repo.GetPaymentSetting(tt.settingID)
				assert.Equal(s.T(), pkgerrors.ErrDataNotFound, err)
			}
		})
	}
}

func TestPaymentSettingsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentSettingsRepositoryTestSuite))
}
