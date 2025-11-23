package handler

import (
	"encoding/json"
	"net/http"
	"pr-reviewer-2025/internal/models"
)

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := models.ErrorResponse{
		Error: models.ErrorBody{
			Code:    code,
			Message: msg,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
