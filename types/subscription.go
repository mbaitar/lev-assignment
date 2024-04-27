package types

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                string    `bun:"id,pk"`
	Customer          string    `bun:"customer"`
	UserID            uuid.UUID `bun:"user_id"`
	Status            string    `bun:"status"`
	CreatedAt         time.Time `bun:"created_at"`
	Amount            float64   `bun:"amount"`
	Currency          string    `bun:"currency"`
	EndDate           time.Time `bun:"end_date"`
	CancelAtPeriodEnd bool      `bun:"cancel_at_period_end"`
}
