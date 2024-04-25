package strp

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/subscription"
	"github.com/uptrace/bun"
	"log/slog"
	"os"
	"time"
)

func FetchAllSubscriptions(token types.IntegrationToken, userID uuid.UUID) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	params := &stripe.SubscriptionListParams{}
	params.SetStripeAccount(token.ConnectedAccountId)

	iter := subscription.List(params)

	for iter.Next() {
		sub := iter.Subscription()

		var totalPrice float64

		for _, subItem := range sub.Items.Data {
			if subItem.Price != nil {
				totalPrice += float64(subItem.Price.UnitAmountDecimal) / 100 * float64(subItem.Quantity)
			}
		}

		var endDate time.Time
		if sub.CancelAt != 0 {
			endDate = time.Unix(sub.CancelAt, 0)
		} else if sub.CurrentPeriodEnd > 0 {
			endDate = time.Unix(sub.CurrentPeriodEnd, 0)
		}

		subscriptionData := &types.Subscription{
			ID:                sub.ID,
			Customer:          sub.Customer.ID,
			UserID:            userID,
			Status:            string(sub.Status),
			CreatedAt:         time.Unix(sub.Created, 0),
			Amount:            totalPrice,
			Currency:          string(sub.Currency),
			EndDate:           endDate,
			CancelAtPeriodEnd: sub.CancelAtPeriodEnd,
		}

		err := db.Bun.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
			if err := db.CreateSubscription(tx, subscriptionData); err != nil {
				slog.Error("token save failed: ", "err", err)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
