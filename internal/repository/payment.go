package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

func (r *repository) CreatePayment(ctx context.Context, orderID string, payment models.Payment) error {
	sql, args, err := sqrl.Insert("payment").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "transaction", "request_id", "currency", "provider",
		"amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee").
		Values(orderID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
			payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).
		ToSql()

	if err != nil {
		return err
	}

	query := database.Query{
		Name:     "CreatePayment",
		QueryRaw: sql,
	}

	r.db.DB().QueryRowContext(ctx, query, args...)
	return nil
}

func (r *repository) GetPaymentByID(ctx context.Context, id string) (models.Payment, error) {
	sql, args, err := sqrl.Select("order_uid", "transaction", "request_id", "currency", "provider",
		"amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee").
		PlaceholderFormat(sqrl.Dollar).From("payment").Where(sqrl.Eq{"order_uid": id}).ToSql()

	if err != nil {
		return models.Payment{}, err
	}

	query := database.Query{
		Name:     "GetPayment",
		QueryRaw: sql,
	}

	var payment models.Payment
	if err := r.db.DB().ScanOneContext(ctx, &payment, query, args...); err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}
