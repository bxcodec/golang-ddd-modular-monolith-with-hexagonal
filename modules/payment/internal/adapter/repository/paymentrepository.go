package repository

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/dbutils"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/uniqueid"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) ports.IPaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) qb() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r *paymentRepository) CreatePayment(p *payment.Payment) (err error) {
	p.ID, err = uniqueid.GeneratePK("pay")
	if err != nil {
		return err
	}

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err = r.qb().Insert("payment_module.payments").
		Columns("id", "amount", "currency", "status", "created_at", "updated_at").
		Values(p.ID, p.Amount, p.Currency, p.Status, p.CreatedAt, p.UpdatedAt).
		RunWith(r.db).
		Exec()
	if err != nil {
		return dbutils.HandlePostgresError(err)
	}

	return nil
}

func (r *paymentRepository) GetPayment(id string) (p payment.Payment, err error) {
	err = r.qb().Select("id", "amount", "currency", "status", "created_at", "updated_at").
		From("payment_module.payments").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&p.ID, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return payment.Payment{}, dbutils.HandlePostgresError(err)
	}

	return p, nil
}

func (r *paymentRepository) FetchPayments(params payment.FetchPaymentsParams) (result []payment.Payment, nextCursor string, err error) {
	query := r.qb().Select("id", "amount", "currency", "status", "created_at", "updated_at").
		From("payment_module.payments").
		OrderBy("id DESC")

	// Apply cursor-based pagination using ULID
	if params.Cursor != "" {
		cursorID, decodeErr := dbutils.DecodeCursor(params.Cursor)
		if decodeErr != nil {
			return nil, "", decodeErr
		}
		query = query.Where(sq.Lt{"id": cursorID})
	}

	if params.Currency != "" {
		query = query.Where(sq.Eq{"currency": params.Currency})
	}

	if params.Status != "" {
		query = query.Where(sq.Eq{"status": params.Status})
	}

	// Fetch one extra to determine if there's a next page
	query = query.Limit(uint64(params.Limit + 1))

	rows, err := query.RunWith(r.db).Query()
	if err != nil {
		return nil, "", err
	}
	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			log.Error().Err(errClose).Msg("failed to close rows")
		}
	}()

	result = make([]payment.Payment, 0)
	for rows.Next() {
		var p payment.Payment
		if err := rows.Scan(&p.ID, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, "", err
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	// Check if there are more results
	if len(result) > params.Limit {
		// Remove the extra item
		result = result[:params.Limit]
		// Set next cursor to the last item's ID
		nextCursor = dbutils.EncodeCursor(result[len(result)-1].ID)
	}

	return result, nextCursor, nil
}

func (r *paymentRepository) UpdatePayment(p *payment.Payment) (err error) {
	p.UpdatedAt = time.Now()

	result, err := r.qb().Update("payment_module.payments").
		Set("amount", p.Amount).
		Set("currency", p.Currency).
		Set("status", p.Status).
		Set("updated_at", p.UpdatedAt).
		Where(sq.Eq{"id": p.ID}).
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
		return errors.ErrDataNotFound
	}

	return nil
}

func (r *paymentRepository) DeletePayment(id string) (err error) {
	result, err := r.qb().Delete("payment_module.payments").
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
		return errors.ErrDataNotFound
	}

	return nil
}
