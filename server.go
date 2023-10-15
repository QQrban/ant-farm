package main

import (
	"encoding/json"
	"lem-in/types"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := types.AllData
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
