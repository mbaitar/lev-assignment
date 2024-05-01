package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/pkg/metrics"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/mbaitar/levenue-assignment/view/dashboard"
)

func HandleDashboardIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	if user.Type == "SELLER" {
		metrics, err := db.GetLatestMetricByUserID(user.ID)
		if err != nil {
			return err
		}
		chartDataMrr := CreateChartData()
		data := dashboard.ViewData{
			Metrics:       metrics,
			MrrTimeValues: chartDataMrr,
		}

		return render(r, w, dashboard.Index(data))
	}
	return render(r, w, dashboard.Index(dashboard.ViewData{}))
}

func CreateChartData() []types.TimeValue {
	// Data to be converted into []TimeValue
	data := []struct {
		Time  string
		Value float64
	}{
		{Time: "2019-04-11", Value: 80.01},
		{Time: "2019-04-12", Value: 96.63},
		{Time: "2019-04-13", Value: 76.64},
		{Time: "2019-04-14", Value: 81.89},
		{Time: "2019-04-15", Value: 74.43},
		{Time: "2019-04-16", Value: 80.01},
		{Time: "2019-04-17", Value: 96.63},
		{Time: "2019-04-18", Value: 76.64},
		{Time: "2019-04-19", Value: 81.89},
		{Time: "2019-04-20", Value: 74.43},
	}

	// Slice to hold the TimeValue structs
	var timeValues []types.TimeValue

	// Parse time strings and create TimeValue structs
	for _, d := range data {
		parsedTime, err := time.Parse("2006-01-02", d.Time)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err)
			continue
		}
		timeValues = append(timeValues, types.TimeValue{
			Time:  parsedTime,
			Value: d.Value,
		})
	}
	return timeValues
}

func HandleTradeCreate(w http.ResponseWriter, r *http.Request) error {
	arrStr := r.FormValue("arr")
	arr, err := strconv.ParseFloat(arrStr, 64)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return err
	}
	user := getAuthenticatedUser(r)
	trade := metrics.CalculateTrade(arr, user.ID)
	err = db.CreateTrade(trade)
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/history")
}

func HandleRunMetricCalculation(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	err := metrics.CalculateMetrics(user.ID)
	if err != nil {
		return err
	}
	metrics, err := db.GetLatestMetricByUserID(user.ID)
	if err != nil {
		return err
	}

	data := dashboard.ViewData{
		Metrics: metrics,
	}

	return render(r, w, dashboard.Index(data))
}
