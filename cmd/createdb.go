package main

import "database/sql"

func initDB(db *sql.DB) error {
	initsql := `
	-- Teams Table
CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR PRIMARY KEY
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR PRIMARY KEY,
    username VARCHAR NOT NULL,
    team_name VARCHAR REFERENCES teams(team_name) ON DELETE SET NULL,
    is_active BOOLEAN 
);

-- PR Table
CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR PRIMARY KEY,
    pull_request_name VARCHAR NOT NULL,
    author_id VARCHAR NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status VARCHAR NOT NULL ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP NULL
);

-- PR Reviewers table 
CREATE TABLE IF NOT EXISTS pr_reviewers (
    pull_request_id VARCHAR NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id VARCHAR NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

-- indexing
CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer ON pr_reviewers(reviewer_id);
	`

	_, err := db.Exec(initsql)
	if err != nil {
		return err
	}
	return nil
}
