package repository

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/dbutils"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/uniqueid"
)

type PaymentSettingsRepository struct {
	db *sql.DB
}

func NewPaymentSettingsRepository(db *sql.DB) (repo *PaymentSettingsRepository) {
	return &PaymentSettingsRepository{
		db: db,
	}
}

func (r *PaymentSettingsRepository) qb() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r *PaymentSettingsRepository) FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (res []paymentsettings.PaymentSetting, nextCursor string, err error) {
	query := r.qb().Select("id", "setting_key", "setting_value", "currency", "status", "created_at", "updated_at").
		From("payment_settings").
		OrderBy("id DESC")

	// Apply currency filter if provided
	if params.Currency != "" {
		query = query.Where(sq.Eq{"currency": params.Currency})
	}

	// Apply setting key filter if provided
	if params.SettingKey != "" {
		query = query.Where(sq.Eq{"setting_key": params.SettingKey})
	}

	// Apply cursor-based pagination using ULID
	if params.Cursor != "" {
		cursorID, decodeErr := dbutils.DecodeCursor(params.Cursor)
		if decodeErr != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", decodeErr)
		}
		query = query.Where(sq.Lt{"id": cursorID})
	}

	// Apply limit
	limit := params.Limit
	if limit <= 0 {
		limit = 10 // default limit
	}
	// Fetch one extra to determine if there's a next page
	query = query.Limit(uint64(limit + 1))

	rows, err := query.RunWith(r.db).Query()
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	res = make([]paymentsettings.PaymentSetting, 0)
	for rows.Next() {
		var setting paymentsettings.PaymentSetting
		if err := rows.Scan(&setting.ID, &setting.SettingKey, &setting.SettingValue, &setting.Currency, &setting.Status, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
			return nil, "", err
		}
		res = append(res, setting)
	}

	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	// Check if there are more results
	if len(res) > limit {
		// Remove the extra item
		res = res[:limit]
		// Set next cursor to the last item's ID
		nextCursor = dbutils.EncodeCursor(res[len(res)-1].ID)
	}

	return res, nextCursor, nil
}

func (r *PaymentSettingsRepository) GetPaymentSettingByCurrency(currency string) (setting paymentsettings.PaymentSetting, err error) {
	err = r.qb().Select("id", "setting_key", "setting_value", "currency", "status", "created_at", "updated_at").
		From("payment_settings").
		Where(sq.Eq{"currency": currency}).
		RunWith(r.db).
		QueryRow().
		Scan(&setting.ID, &setting.SettingKey, &setting.SettingValue, &setting.Currency, &setting.Status, &setting.CreatedAt, &setting.UpdatedAt)

	if err != nil {
		return paymentsettings.PaymentSetting{}, dbutils.HandlePostgresError(err)
	}

	return setting, nil
}

func (r *PaymentSettingsRepository) CreatePaymentSetting(settings *paymentsettings.PaymentSetting) (err error) {
	settings.ID, err = uniqueid.GeneratePK("pset")
	if err != nil {
		return err
	}

	now := time.Now()
	settings.CreatedAt = now
	settings.UpdatedAt = now

	_, err = r.qb().Insert("payment_settings").
		Columns("id", "setting_key", "setting_value", "currency", "status", "created_at", "updated_at").
		Values(settings.ID, settings.SettingKey, settings.SettingValue, settings.Currency, settings.Status, settings.CreatedAt, settings.UpdatedAt).
		RunWith(r.db).
		Exec()

	if err != nil {
		return dbutils.HandlePostgresError(err)
	}

	return nil
}

func (r *PaymentSettingsRepository) UpdatePaymentSetting(settings *paymentsettings.PaymentSetting) (err error) {
	settings.UpdatedAt = time.Now()

	result, err := r.qb().Update("payment_settings").
		Set("setting_key", settings.SettingKey).
		Set("setting_value", settings.SettingValue).
		Set("currency", settings.Currency).
		Set("status", settings.Status).
		Set("updated_at", settings.UpdatedAt).
		Where(sq.Eq{"id": settings.ID}).
		RunWith(r.db).
		Exec()

	if err != nil {
		return dbutils.HandlePostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return dbutils.HandlePostgresError(err)
	}

	return nil
}

func (r *PaymentSettingsRepository) DeletePaymentSetting(id string) (err error) {
	result, err := r.qb().Delete("payment_settings").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Exec()

	if err != nil {
		return dbutils.HandlePostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return dbutils.HandlePostgresError(err)
	}

	return nil
}
