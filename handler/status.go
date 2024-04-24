package handler

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	App     string
	Version string
	Status  string
}

func HandleStatus(w http.ResponseWriter, r *http.Request) error {
	s := &Status{
		App:     "lev-assignment",
		Version: "v1",
		Status:  "Operational",
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(s)
}
