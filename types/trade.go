package types

import (
	"time"

	"github.com/google/uuid"
)

type Trade struct {
	TradeID       int       `bun:"trade_id,pk,autoincrement"`
	Buyer         uuid.UUID `bun:"buyer"`
	ARRTraded     float64   `bun:"arr_traded"`
	DiscountRate  float64   `bun:"discount_rate"`
	TradeDate     time.Time `bun:"trade_date,default:current_timestamp"`
	ROI           float64   `bun:"roi"`
	ROIPercentage float64   `bun:"roi_percentage"`
	NetProfit     float64   `bun:"net_profit"`
}
