package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"pr-reviewer-2025/internal/models"
	"pr-reviewer-2025/internal/repository"
	"strings"
)

func (h *Handler) createTeam(w http.ResponseWriter, r *http.Request) {
	var req models.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID", "invalid json")
		return
	}
	defer r.Body.Close()

	team, err := h.services.Teams.CreateOrUpdateTeam(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrTeamExists) {
			h.writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]models.Team{"team": team})
}

func (h *Handler) getTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		teamName = chi.URLParam(r, "team_name")
	}

	if teamName == "" {
		h.writeError(w, http.StatusBadRequest, "INVALID", "team_name required")
		return
	}

	team, err := h.services.Teams.GetTeam(r.Context(), teamName)
	if err != nil {
		lerr := strings.ToLower(err.Error())
		if strings.Contains(lerr, "not found") || strings.Contains(lerr, "no rows") {
			h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, team)
}
