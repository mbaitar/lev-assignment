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
	mrr, err := db.CalculateMRR(user)
	if err != nil {
		return err
	}

	churnAmount, err := db.CalculateChurn(from, user)
	if err != nil {
		return err
	}

	churnMRR, err := db.CalculateChurnedMRR(from, user)
	if err != nil {
		return err
	}

	churnPercentage, err := db.CalculateChurnPercentage(from, user)
	if err != nil {
		return err
	}

	netGrowth, err := db.CalculateNetGrowth(from, user)
	if err != nil {
		return err
	}

	tradingLimit := mrr * 30 / 100
	metric := &types.Metric{
		User:              user,
		MRR:               mrr,
		ChurnAmount:       churnAmount,
		ChurnedMRR:        churnMRR,
		ChurnedPercentage: churnPercentage,
		NetGrowth:         int64(netGrowth),
		TradingLimit:      tradingLimit,
	}
	err = db.CreateMetrics(metric)
	if err != nil {
		slog.Error("metric save error", "err", err)
		return err
	}
	return nil
}
