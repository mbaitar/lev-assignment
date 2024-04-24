package types

import (
	"github.com/google/uuid"
	"time"
)

type IntegrationToken struct {
	TokenID            int       `bun:"token_id,pk,autoincrement"`
	AccessToken        string    `bun:"access_token,notnull"`
	RefreshToken       string    `bun:"refresh_token,notnull"`
	ConnectedAccountId string    `bun:"connected_account_id,notnull"`
	CreatedAt          time.Time `bun:"created_at,default:current_timestamp"`
	AccountID          uuid.UUID `bun:"account_id,unique,notnull"`
}
