package models

import "time"

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

type PR struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type PRsmall struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type AnswerbyReassignRequest struct {
	PR    *PR    `json:"pull_request"`
	NewId string `json:"replaced_by"`
}

type GeneralStats struct {
	TotalOpen   int `json:"total_open"`
	TotalMerged int `json:"total_merged"`
}

type UserWorkload struct {
	UserName        string `json:"user_name"`
	CreatedTotal    int    `json:"created_total"`
	AssignedPending int    `json:"assigned_pending"`
}

type StatsResponnse struct {
	GlobalStat GeneralStats    `json:"Global_stats"`
	UserStat   []*UserWorkload `json:"User_stats"`
}
