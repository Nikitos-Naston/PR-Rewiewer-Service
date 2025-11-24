package repository

import (
	"PRreviewService/internal/models"
	"database/sql"
	"fmt"
)

var (
	tableTeam string = "teams"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (tr *TeamRepository) CreateTeam(teamName string) error {
	query := fmt.Sprintf("INSERT INTO %s(team_name) VALUES($1)", tableTeam)
	if _, err := tr.db.Query(query, teamName); err != nil {
		return err

	}
	return nil
}

func (tr *TeamRepository) TeamExist(teamName string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE team_name = $1 LIMIT 1", tableTeam)
	var a int

	err := tr.db.QueryRow(query, teamName).Scan(&a)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (tr *TeamRepository) GetTeam(teamName string) (*models.Team, error) {
	team := models.Team{}

	exist, err := tr.TeamExist(teamName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, err
	}

	userRepo := NewUserRepository(tr.db)
	users, err := userRepo.GetAllByTeam(teamName)
	if err != nil {
		return nil, err
	}
	team.Members = users
	return &team, nil
}
