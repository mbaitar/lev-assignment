package strp

import (
	"context"
	"database/sql"
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/subscription"
	"github.com/uptrace/bun"
	"log/slog"
	"os"
	"time"
)

func FetchAllSubscriptions(token types.IntegrationToken) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY") // My api key

	// Fetch all subscriptions from Stripe
	params := &stripe.SubscriptionListParams{}
	params.Filters.AddFilter("limit", "", "100") // Fetch first 100 subscriptions
	params.SetStripeAccount(token.ConnectedAccountId)

	iter := subscription.List(params)

	for iter.Next() {
		sub := iter.Subscription()

		var totalPrice int64

		for _, subItem := range sub.Items.Data {
			if subItem.Price != nil {
				totalPrice += subItem.Price.UnitAmount * subItem.Quantity
			}
		}

		subscriptionData := &types.Subscription{
			ID:        sub.ID,
			Customer:  sub.Customer.ID,
			Status:    string(sub.Status),
			CreatedAt: time.Unix(sub.Created, 0),
			Amount:    int(totalPrice),
			Currency:  string(sub.Currency),
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
