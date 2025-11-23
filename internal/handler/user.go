package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"pr-reviewer-2025/internal/models"
	"strings"
)

func (h *Handler) setIsActive(w http.ResponseWriter, r *http.Request) {
	var req models.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID", "invalid json")
		return
	}
	defer r.Body.Close()

	user, err := h.services.Users.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		le := strings.ToLower(err.Error())
		if strings.Contains(le, "not found") || strings.Contains(le, "user not found") {
			h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]models.User{"user": user})
}

func (h *Handler) getUserReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = chi.URLParam(r, "user_id")
	}
	if userID == "" {
		h.writeError(w, http.StatusBadRequest, "INVALID", "user_id required")
		return
	}

	prs, err := h.services.Users.GetPRsForReviewer(r.Context(), userID)
	if err != nil {
		le := strings.ToLower(err.Error())
		if strings.Contains(le, "not found") || strings.Contains(le, "no rows") {
			h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		return
	}

	resp := map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	}
	h.writeJSON(w, http.StatusOK, resp)
}
