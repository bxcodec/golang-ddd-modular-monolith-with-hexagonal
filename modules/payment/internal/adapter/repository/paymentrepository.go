package repository

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/dbutils"
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

	_, err = r.qb().Insert("payments").
		Columns("id", "amount", "currency", "status", "created_at", "updated_at").
		Values(p.ID, p.Amount, p.Currency, p.Status, p.CreatedAt, p.UpdatedAt).
		RunWith(r.db).
		Exec()

	if err != nil {
		return dbutils.HandlePostgresError(err)
	}

	return nil
}

func (r *paymentRepository) GetPayment(id string) (p *payment.Payment, err error) {
	var paymentData payment.Payment

	err = r.qb().Select("id", "amount", "currency", "status", "created_at", "updated_at").
		From("payments").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&paymentData.ID, &paymentData.Amount, &paymentData.Currency, &paymentData.Status, &paymentData.CreatedAt, &paymentData.UpdatedAt)

	if err != nil {
		return nil, dbutils.HandlePostgresError(err)
	}

	return &paymentData, nil
}

func (r *paymentRepository) GetPayments() (payments []*payment.Payment, err error) {
	rows, err := r.qb().Select("id", "amount", "currency", "status", "created_at", "updated_at").
		From("payments").
		OrderBy("created_at DESC").
		RunWith(r.db).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments = make([]*payment.Payment, 0)
	for rows.Next() {
		var p payment.Payment
		if err := rows.Scan(&p.ID, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentRepository) UpdatePayment(p *payment.Payment) (err error) {
	p.UpdatedAt = time.Now()

	result, err := r.qb().Update("payments").
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
		return dbutils.HandlePostgresError(err)
	}

	return nil
}

func (r *paymentRepository) DeletePayment(id string) (err error) {
	result, err := r.qb().Delete("payments").
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
