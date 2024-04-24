package types

import (
	"github.com/google/uuid"
	"time"
)

type Metric struct {
	ID           int       `bun:"id,pk,autoincrement"`
	User         uuid.UUID `bun:"user_id,notnull"`
	MRR          float64   `bun:"mrr,notnull"`
	Churn        float64   `bun:"churn,notnull"`
	NetGrowth    int64     `bun:"net_growth,notnull"`
	TradingLimit float64   `bun:"trading_limit,notnull"`
	LastUpdated  time.Time `bun:"last_updated"`
}
