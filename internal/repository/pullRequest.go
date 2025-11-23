package repository

import (
	"PRreviewService/internal/models"
	"database/sql"
	"fmt"
)

var (
	tablePR string = "pull_requests"
	tablePRreviewer string = "pr_reviewers"
)

type PRRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (prr *PRRepository) CreatePR(pullRequest models.PR) models.PR, error {
	query := fmt.Sprintf("INSERT INTO %s (pull_request_id, pull_request_name, author_id, status, created_at)) VALUES ($1, $2, $3, 'OPEN', CURRENT_TIMESTAMP)", tablePR)
	if _, err := ppr.db.Query(query, pull_requests.ID, pull_requests.Name, pull_requests.AuthorID); err != nil {
		return nil, err

	}
	return pull_requests, nil
}

func (prr *PRRepository) GetPRByID (pullRequestID string) (models.PR, error) {
	query := fmt.Sprintf("SELECT (pull_request_id, pull_request_name, author_id, status, created_at, merged_at) FROM %s WHERE pull_request_id = $1", tablePR)
	var pr models.PR
	err := prr.db.QueryRow(query, pullRequestID).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.created_at, &pr.merged_at)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (prr *PRRepository) GetReviewers (pullRequestID string) ([]string, error) {
	query := fmt.Sprintf("SELECT reviewer_id FROM %s WHERE pull_request_id = $1 ORDER BY reviewer_id", tablePRreviewer)
	rows, err := ur.db.Query(query, pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var rewiewerID string
		err := rows.Scan(&rewiewerID)
		if err != nil {
			messages.SendLogMessage(logrus.ErrorLevel, "Cant parse rewiewr", err)
		}
		reviewers = append(reviewers, rewiewerID)
	}
	return reviewers, nil
}
