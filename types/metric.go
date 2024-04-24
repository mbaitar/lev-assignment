package types

import (
	"github.com/google/uuid"
	"time"
)

type Metric struct {
	ID           int       `bun:"id,pk,autoincrement"`
	User         uuid.UUID `bun:"user_id,notnull"`
	MRR          int64     `bun:"mrr,notnull"`
	Churn        int       `bun:"churn,notnull"`
	NetGrowth    int       `bun:"net_growth,notnull"`
	TradingLimit int64     `bun:"trading_limit,notnull"`
	LastUpdated  time.Time `bun:"last_updated"`
}
