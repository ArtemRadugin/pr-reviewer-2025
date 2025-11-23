package models

import (
	"github.com/lib/pq"
	"time"
)

type PullRequest struct {
	PullRequestID     string         `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName   string         `json:"pull_request_name" db:"pull_request_name"`
	AuthorID          string         `json:"author_id" db:"author_id"`
	Status            string         `json:"status" db:"status"` // enum: OPEN, MERGED
	AssignedReviewers pq.StringArray `db:"assigned_reviewers" json:"assigned_reviewers"`
	CreatedAt         *time.Time     `json:"createdAt,omitempty" db:"created_at"`
	MergedAt          *time.Time     `json:"mergedAt,omitempty" db:"merged_at"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string `json:"author_id" db:"author_id"`
	Status          string `json:"status" db:"status"` // enum: OPEN, MERGED
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePRResponse struct {
	PR PullRequest `json:"pr"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePRResponse struct {
	PR PullRequest `json:"pr"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}
