package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer-2025/internal/models"
	"pr-reviewer-2025/internal/repository"
	"strings"
)

func (h *Handler) createPR(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID", "invalid json")
		return
	}

	defer r.Body.Close()

	pr, err := h.services.PullRequests.CreatePR(r.Context(), req)
	if err != nil {
		lerr := strings.ToLower(err.Error())
		switch {
		case strings.Contains(lerr, "author not found") || strings.Contains(lerr, "team not found"):
			h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		case strings.Contains(lerr, "pr exists") || strings.Contains(lerr, "already exists"):
			h.writeError(w, http.StatusConflict, "PR_EXISTS", err.Error())
			return
		default:
			h.writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
			return
		}
	}

	h.writeJSON(w, http.StatusCreated, map[string]models.PullRequest{"pr": pr})
}

func (h *Handler) mergePR(w http.ResponseWriter, r *http.Request) {
	var req models.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID", "invalid json")
		return
	}
	defer r.Body.Close()

	pr, err := h.services.PullRequests.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		lerr := strings.ToLower(err.Error())
		if strings.Contains(lerr, "not found") || strings.Contains(lerr, "pr not found") {
			h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]models.PullRequest{"pr": pr})
}

func (h *Handler) reassign(w http.ResponseWriter, r *http.Request) {
	var req models.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID", "invalid json")
		return
	}
	defer r.Body.Close()

	if req.PullRequestID == "" || req.OldReviewerID == "" {
		h.writeError(w, http.StatusBadRequest, "INVALID", "pull_request_id and old_reviewer_id required")
		return
	}

	pr, replacedBy, err := h.services.PullRequests.ReassignReviewer(
		r.Context(),
		req.PullRequestID,
		req.OldReviewerID,
	)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPRNotFound):
			h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		case errors.Is(err, repository.ErrPRMerged):
			h.writeError(w, http.StatusConflict, "PR_MERGED", err.Error())
			return
		case errors.Is(err, repository.ErrNotAssigned):
			h.writeError(w, http.StatusConflict, "NOT_ASSIGNED", err.Error())
			return
		case errors.Is(err, repository.ErrNoCandidate):
			h.writeError(w, http.StatusConflict, "NO_CANDIDATE", err.Error())
			return
		default:
			h.writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
			return
		}
	}

	resp := map[string]interface{}{
		"pr":          pr,
		"replaced_by": replacedBy,
	}

	h.writeJSON(w, http.StatusOK, resp)
}
