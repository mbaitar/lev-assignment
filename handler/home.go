package handler

import (
	"github.com/mbaitar/levenue-assignment/view/home"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, home.Index())
}
