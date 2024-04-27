package handler

import (
	"net/http"

	"github.com/mbaitar/levenue-assignment/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, home.Index())
}
