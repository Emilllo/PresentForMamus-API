package handlers

import (
	"net/http"
	"strconv"
)

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return 0, false
	}

	return id, true
}

func parseIDFromPath(w http.ResponseWriter, r *http.Request, param string) (int, bool) {
	id, err := strconv.Atoi(r.PathValue(param))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid "+param, http.StatusBadRequest)
		return 0, false
	}

	return id, true
}
