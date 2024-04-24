package metrics

import (
	"github.com/google/uuid"
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/types"
	"log/slog"
	"time"
)

func CalculateMetrics(user uuid.UUID) error {
	now := time.Now()
	from := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	mrr, err := db.CalculateMRR()
	if err != nil {
		return err
	}
	churn, err := db.CalculateChurn(from)
	if err != nil {
		return err
	}
	netGrowth, err := db.CalculateNetGrowth(from)
	if err != nil {
		return err
	}
	tradingLimit := db.CalculateTradingLimit(mrr)
	if err != nil {
		return err
	}
	metric := &types.Metric{
		User:         user,
		MRR:          mrr,
		Churn:        churn,
		NetGrowth:    netGrowth,
		TradingLimit: tradingLimit,
	}
	err = db.CreateMetrics(metric)
	if err != nil {
		slog.Error("metric save error", "err", err)
		return err
	}
	return nil
}
