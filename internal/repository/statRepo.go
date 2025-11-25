package repository

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"database/sql"
	"fmt"
)

// used dependencies
// var (
// 	tableTeam       string = "teams"
// 	tableUser       string = "users"
// 	tablePR         string = "pull_requests"
// 	tablePRreviewer string = "pr_reviewers"
// )

type StatRepository struct {
	db *sql.DB
}

func NewStatRepository(db *sql.DB) *StatRepository {
	return &StatRepository{db: db}
}

func (sr *StatRepository) GetPRStats() (*models.GeneralStats, error) {
	query := fmt.Sprintf(`
        SELECT 
            COUNT(*) FILTER (WHERE status = 'OPEN') as open_count,
            COUNT(*) FILTER (WHERE status = 'MERGED') as merged_count
        FROM %s
    `, tablePR)
	var stats models.GeneralStats

	err := sr.db.QueryRow(query).Scan(&stats.TotalOpen, &stats.TotalMerged)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (sr *StatRepository) GetUserStats() ([]*models.UserWorkload, error) {
	query := fmt.Sprintf(`
        SELECT 
            u.username,
            SUM(created_cnt) as created_total,
            SUM(review_cnt) as assigned_pending
        FROM (
            SELECT author_id as user_id, COUNT(*) as created_cnt, 0 as review_cnt 
            FROM %s 
            GROUP BY author_id
            
            UNION ALL
        
            SELECT 
                r.reviewer_id as user_id, 
                0 as created_cnt, 
                COUNT(*) as review_cnt 
            FROM pr_reviewers r
            JOIN %s pr ON r.pull_request_id = pr.pull_request_id
            WHERE pr.status = 'OPEN' 
            GROUP BY r.reviewer_id
        ) combined_stats
        JOIN users u ON u.user_id = combined_stats.user_id
        WHERE combined_stats.user_id IS NOT NULL
        GROUP BY u.username
        ORDER BY assigned_pending DESC, created_total DESC;
    `, tablePR, tablePR)

	rows, err := sr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.UserWorkload, 0)
	for rows.Next() {
		var u models.UserWorkload
		if err := rows.Scan(&u.UserName, &u.CreatedTotal, &u.AssignedPending); err != nil {
			messages.SendLogMessage("Cant parse user", err)
		}
		result = append(result, &u)
	}

	return result, nil

}
