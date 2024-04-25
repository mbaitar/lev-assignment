package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/subscription"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const numWorkers = 10

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	stripe.Key = os.Getenv("STRIPE_API_KEY_TEST_USER")

	// var priceIDs = []string{"price_1P8OEWD9fA4p74AFfD1tXrlH", "price_1P8OErD9fA4p74AFhtijxfQn", "price_1P8OF6D9fA4p74AFKXwpcz37"} // 9.99 - 19.99 - 49.99
	var priceIDs = []string{"price_1P901hK6VtaRtZVAn2fNTWhP", "price_1P901UK6VtaRtZVAdLdmbrq5", "price_1P901LK6VtaRtZVAm4dveHMe"}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for workerID := 1; workerID <= numWorkers; workerID++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				generateDummyData(rng, priceIDs, workerID)
			}
		}(workerID)
	}

	wg.Wait()

	slog.Info("Completed creating 1,000 customers and subscriptions")
}

func generateDummyData(rng *rand.Rand, priceIDs []string, workerID int) {
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(faker.Email()),
		Name:  stripe.String(faker.Name()),
	}

	newCustomer, err := customer.New(customerParams)
	if err != nil {
		slog.Error("Creating customer", "Worker", workerID, "err", err)
		return
	}

	priceID := priceIDs[rng.Intn(len(priceIDs))]

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(newCustomer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		CollectionMethod: stripe.String(string(stripe.SubscriptionCollectionMethodSendInvoice)),
		DaysUntilDue:     stripe.Int64(30),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			PaymentMethodTypes: []*string{
				stripe.String(string(stripe.SubscriptionPaymentSettingsPaymentMethodTypeCustomerBalance)),
			},
		},
	}

	churn := rng.Float64() < 0.15 // 15% churn
	if churn {
		subscriptionParams.CancelAtPeriodEnd = stripe.Bool(true)
	}

	_, err = subscription.New(subscriptionParams)
	if err != nil {
		slog.Error("Creating subscription", "Worker", workerID, "err", err)
		return
	}

	slog.Info("Created customer and subscription", "Worker", workerID)
}

func initEverything() error {
	return godotenv.Load()
}
