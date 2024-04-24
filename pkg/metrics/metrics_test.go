package metrics

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Ensure SQLite is imported
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateTrade(t *testing.T) {
	discountRate := 0.88
	arr := 500000.0
	userID := uuid.New()

	// Expected results
	investment := arr * discountRate
	roi := arr - investment
	expectedROIPercentage := (roi / investment) * 100

	trade := CalculateTrade(arr, userID)

	assert.NotNil(t, trade, "Trade should not be nil")
	assert.Equal(t, userID, trade.Buyer, "Wrong buyer ID got=%f, wanted=%f", trade.Buyer, userID)
	assert.Equal(t, arr, trade.ARRTraded, "Wrong ArrTraded got=%f, wanted=%f", trade.ARRTraded, arr)
	assert.Equal(t, discountRate, trade.DiscountRate, "Wrong DiscountRate got=%f, wanted=%f", trade.DiscountRate, discountRate)
	assert.Equal(t, roi, trade.NetProfit, "Wrong netProfit got=%f, wanted=%f", trade.NetProfit, roi)
	assert.Equal(t, roi, trade.ROI, "Wrong ROI got=%f, wanted=%f", trade.ROI, roi)
	assert.InDelta(t, expectedROIPercentage, trade.ROIPercentage, 0.01, "Wrong ROIPercentage got=%f, wanted=%f", trade.ROIPercentage, expectedROIPercentage)
}

func setupInMemorySQLiteDB(t *testing.T) *bun.DB {
	sqlDB, err := sql.Open("sqlite3", ":memory:") // Connect to in-memory SQLite
	require.NoError(t, err, "Failed to open SQLite connection")
	bunDB := bun.NewDB(sqlDB, sqlitedialect.New())
	return bunDB
}

func TestCalculateMetrics(t *testing.T) {
	sqlDB := setupInMemorySQLiteDB(t)
	defer sqlDB.Close()

	_, err := sqlDB.NewCreateTable().
		Model((*types.Subscription)(nil)).
		IfNotExists().
		Exec(context.Background())
	require.NoError(t, err, "Failed to create subscriptions table")

	_, err = sqlDB.NewCreateTable().
		Model((*types.Metric)(nil)).
		IfNotExists().
		Exec(context.Background())
	require.NoError(t, err, "Failed to create metrics table")

	testUserID := uuid.New()
	now := time.Now()

	_, err = sqlDB.NewInsert().
		Model(&types.Subscription{
			ID:                "sub_123",
			Customer:          "cus_123",
			UserID:            testUserID,
			Status:            "active",
			CreatedAt:         now,
			Amount:            100.0,
			Currency:          "USD",
			EndDate:           now.AddDate(0, 1, 0),
			CancelAtPeriodEnd: false,
		}).
		Exec(context.Background())
	require.NoError(t, err, "Failed to insert active subscription")

	_, err = sqlDB.NewInsert().
		Model(&types.Subscription{
			ID:                "sub_456",
			Customer:          "cus_456",
			UserID:            testUserID,
			Status:            "canceled",
			CreatedAt:         now.AddDate(0, -1, 0),
			Amount:            200.0,
			Currency:          "USD",
			EndDate:           now.AddDate(0, 1, 2),
			CancelAtPeriodEnd: true,
		}).
		Exec(context.Background())
	require.NoError(t, err, "Failed to insert churned subscription")

	db.Bun = sqlDB

	err = CalculateMetrics(testUserID)
	require.NoError(t, err, "CalculateMetrics failed")

	var metric types.Metric
	err = sqlDB.NewSelect().
		Model(&metric).
		Where("user_id = ?", testUserID).
		Scan(context.Background())
	require.NoError(t, err, "Failed to retrieve metric")

	expectedMrr := 100.0
	expectedChurnedAmount := 1
	expectedChurnedMRR := 200.0
	expectedChurnedPercentage := (expectedChurnedMRR / expectedMrr) * 100
	expectedTradingLimit := float64(30)

	assert.Equal(t, expectedMrr, metric.MRR, "Wrong MRR value got=%f, wanted=%f", metric.MRR, expectedMrr)
	assert.Equal(t, expectedChurnedAmount, metric.ChurnAmount, "Wrong churned amount got=%d, wanted=%d", metric.ChurnAmount, expectedChurnedAmount)
	assert.Equal(t, expectedChurnedMRR, metric.ChurnedMRR, "Wrong churned MRR got=%f, wanted=%f", metric.ChurnedMRR, expectedChurnedMRR)
	assert.Equal(t, expectedChurnedPercentage, metric.ChurnedPercentage, "Wrong churn percentage got=%f, wanted=%f", metric.ChurnedPercentage, expectedChurnedPercentage)
	assert.Equal(t, expectedTradingLimit, metric.TradingLimit, "Wrong trading limit got=%f, wanted=%f", metric.TradingLimit, expectedTradingLimit)
}
