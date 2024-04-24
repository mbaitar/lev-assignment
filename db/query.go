package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/uptrace/bun"
	"time"
)

func GetAccountByUserID(userID uuid.UUID) (types.Account, error) {
	var account types.Account
	err := Bun.NewSelect().
		Model(&account).
		Where("user_id = ?", userID).
		Scan(context.Background())
	return account, err
}

func CreateAccount(account *types.Account) error {
	_, err := Bun.NewInsert().
		Model(account).
		Exec(context.Background())
	return err
}

func CreateIntegrationTokens(token *types.IntegrationToken) error {
	_, err := Bun.NewInsert().
		Model(token).
		On("CONFLICT (account_id) DO UPDATE").
		Set("access_token = EXCLUDED.access_token").
		Set("refresh_token = EXCLUDED.refresh_token").
		Exec(context.Background())
	return err
}

func GetTokenByUserID(userID uuid.UUID) (types.IntegrationToken, error) {
	var token types.IntegrationToken
	err := Bun.NewSelect().
		Model(&token).
		Where("account_id = ?", userID).
		Scan(context.Background())
	return token, err
}

func CreateSubscription(tx bun.Tx, sub *types.Subscription) error {
	_, err := tx.NewInsert().
		Model(sub).
		On("CONFLICT (id) DO UPDATE").
		Set("customer = EXCLUDED.customer").
		Set("status = EXCLUDED.status").
		Set("amount = EXCLUDED.amount").
		Set("currency = EXCLUDED.currency").
		Exec(context.Background())
	return err
}

func CalculateMRR() (int64, error) {
	var totalMRR int64
	err := Bun.NewSelect().
		Model((*types.Subscription)(nil)).
		ColumnExpr("SUM(amount) AS total_mrr").
		Where("status = 'active'").
		Scan(context.Background(), &totalMRR)
	return totalMRR, err
}

func CalculateChurn(fromDate time.Time) (int, error) {
	count, err := Bun.NewSelect().
		Model((*types.Subscription)(nil)).
		Where("status = 'canceled'").
		Where("created_at >= ?", fromDate).
		Count(context.Background())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CalculateNetGrowth(fromDate time.Time) (int, error) {
	newSubs, err := Bun.NewSelect().
		Model((*types.Subscription)(nil)).
		Where("status = 'active'").
		Where("created_at >= ?", fromDate).
		Count(context.Background())

	canceledSubs, err := Bun.NewSelect().
		Model((*types.Subscription)(nil)).
		Where("status = 'canceled'").
		Where("created_at >= ?", fromDate).
		Count(context.Background())

	netGrowth := newSubs - canceledSubs
	return netGrowth, err
}

func CalculateTradingLimit(mrr int64) int64 {
	return mrr * 30 / 100
}

func CreateMetrics(metrics *types.Metric) error {
	_, err := Bun.NewInsert().
		Model(metrics).
		On("CONFLICT (id) DO UPDATE").
		Set("user_id = EXCLUDED.user_id").
		Set("mrr = EXCLUDED.mrr").
		Set("churn = EXCLUDED.churn").
		Set("net_growth = EXCLUDED.net_growth").
		Set("trading_limit = EXCLUDED.trading_limit").
		Set("last_updated = EXCLUDED.last_updated").
		Exec(context.Background())
	return err
}

func GetMetricByUserID(userID uuid.UUID) (types.Metric, error) {
	var metric types.Metric
	err := Bun.NewSelect().
		Model(&metric).
		Where("user_id = ?", userID.String()).
		Scan(context.Background())
	return metric, err
}
