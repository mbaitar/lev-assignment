package types

import "time"

type Subscription struct {
	ID        string    `bun:"id,pk"`
	Customer  string    `bun:"customer"`
	Status    string    `bun:"status"`
	CreatedAt time.Time `bun:"created_at"`
	Amount    int       `bun:"amount"`
	Currency  string    `bun:"currency"`
}
