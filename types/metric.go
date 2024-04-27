package types

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	ID                int       `bun:"id,pk,autoincrement"`
	User              uuid.UUID `bun:"user_id,notnull"`
	MRR               float64   `bun:"mrr,notnull"`
	ChurnAmount       int       `bun:"churned_amount,notnull"`
	ChurnedMRR        float64   `bun:"churned_mrr,notnull"`
	ChurnedPercentage float64   `bun:"churned_percentage,notnull"`
	NetGrowth         int64     `bun:"net_growth,notnull"`
	TradingLimit      float64   `bun:"trading_limit,notnull"`
	LastUpdated       time.Time `bun:"last_updated"`
}
