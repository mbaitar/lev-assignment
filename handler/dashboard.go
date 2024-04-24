package handler

import (
	"fmt"
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/mbaitar/levenue-assignment/view/dashboard"
	"net/http"
	"strconv"
)

const discountRate = 0.88

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
	fmt.Println("hello")
	arrStr := r.FormValue("arr")
	arr, err := strconv.ParseFloat(arrStr, 64)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return err
	}

	investment := arr * discountRate
	netProfit := arr - investment
	roiPercentage := (netProfit / investment) * 100
	user := getAuthenticatedUser(r)
	trade := &types.Trade{
		Buyer:         user.ID,
		ARRTraded:     arr,
		DiscountRate:  discountRate,
		ROI:           netProfit,
		ROIPercentage: roiPercentage,
		NetProfit:     netProfit,
	}
	err = db.CreateTrade(trade)
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/history")
}
