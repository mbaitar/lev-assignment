package handler

import (
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/view/dashboard"
	"net/http"
)

func HandleDashboardIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	metrics, err := db.GetMetricByUserID(user.ID)
	if err != nil {
		return err
	}

	data := dashboard.ViewData{
		Metrics: metrics,
	}

	return render(r, w, dashboard.Index(data))
}
