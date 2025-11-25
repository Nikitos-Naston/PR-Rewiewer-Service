package repository

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"database/sql"
	"fmt"
)

var (
	tablePR         string = "pull_requests"
	tablePRreviewer string = "pr_reviewers"
)

type PRRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) *PRRepository {
	return &PRRepository{db: db}
}

func (prr *PRRepository) CreatePR(pullRequest models.PR) (*models.PR, error) {
	query := fmt.Sprintf("INSERT INTO %s (pull_request_id, pull_request_name, author_id, status, created_at) VALUES ($1, $2, $3, 'OPEN', CURRENT_TIMESTAMP)", tablePR)
	if _, err := prr.db.Exec(query, pullRequest.ID, pullRequest.Name, pullRequest.AuthorID); err != nil {
		return nil, err

	}
	return &pullRequest, nil
}

func (prr *PRRepository) GetPRByID(pullRequestID string) (*models.PR, error) {
	query := fmt.Sprintf("SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at FROM %s WHERE pull_request_id = $1", tablePR)
	var pr models.PR
	err := prr.db.QueryRow(query, pullRequestID).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (prr *PRRepository) GetReviewers(pullRequestID string) ([]string, error) {
	query := fmt.Sprintf("SELECT reviewer_id FROM %s WHERE pull_request_id = $1 ORDER BY reviewer_id", tablePRreviewer)
	rows, err := prr.db.Query(query, pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var rewiewerID string
		err := rows.Scan(&rewiewerID)
		if err != nil {
			messages.SendLogMessage("Cant parse rewiewr", err)
		}
		reviewers = append(reviewers, rewiewerID)
	}
	return reviewers, nil
}

func (prr *PRRepository) CreatePRRewie(pullRequestID string, reviewerID string) error {
	query := fmt.Sprintf("INSERT INTO %s (pull_request_id, reviewer_id) VALUES ($1, $2)", tablePRreviewer)
	if _, err := prr.db.Exec(query, pullRequestID, reviewerID); err != nil {
		return err

	}
	return nil
}

func (prr *PRRepository) DeletePR(pullRequestID string, reviewerID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE pull_request_id = $1 AND reviewer_id = $2", tablePRreviewer)
	if _, err := prr.db.Exec(query, pullRequestID, reviewerID); err != nil {
		return err

	}
	return nil
}

func (prr *PRRepository) MERGEPR(pullRequestID string) error {
	query := fmt.Sprintf("UPDATE %s SET status = 'MERGED', merged_at = CURRENT_TIMESTAMP WHERE pull_request_id = $1", tablePR)
	if _, err := prr.db.Exec(query, pullRequestID); err != nil {
		return err

	}
	return nil
}

func (prr *PRRepository) IsReviewerAssigned(pullRequestID string, reviewerID string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE pull_request_id = $1 AND reviewer_id = $2 LIMIT 1", tablePRreviewer)
	var a int

	err := prr.db.QueryRow(query, pullRequestID, reviewerID).Scan(&a)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (prr *PRRepository) PRExist(pullRequestID string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE pull_request_id = $1  LIMIT 1", tablePR)
	var a int

	err := prr.db.QueryRow(query, pullRequestID).Scan(&a)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (prr *PRRepository) GetRRByUserID(UserID string) ([]string, error) {
	query := fmt.Sprintf("SELECT pull_request_id FROM %s WHERE reviewer_id = $1 ORDER BY pull_request_id", tablePRreviewer)
	rows, err := prr.db.Query(query, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]string, 0)
	for rows.Next() {
		var requstID string
		err := rows.Scan(&requstID)
		if err != nil {
			messages.SendLogMessage("Cant parse rewiewr", err)
		}
		requests = append(requests, requstID)
	}
	return requests, nil

}
