package handler

import (
	"net/http"

	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/view/history"
)

func HandleHistoryIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	trades, err := db.GetTradesByUserID(user.ID)
	if err != nil {
		return err
	}

	data := history.HistoryData{
		Trades: trades,
	}

	return render(r, w, history.Index(data))
}
