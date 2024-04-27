package handler

import (
	"net/http"
	"strconv"

	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/pkg/metrics"
	"github.com/mbaitar/levenue-assignment/view/dashboard"
)

func HandleDashboardIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	if user.Type == "SELLER" {
		metrics, err := db.GetMetricByUserID(user.ID)
		if err != nil {
			return err
		}

		data := dashboard.ViewData{
			Metrics: metrics,
		}

		return render(r, w, dashboard.Index(data))
	}
	return render(r, w, dashboard.Index(dashboard.ViewData{}))
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
