package repository

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"database/sql"
	"fmt"
)

var (
	tableUser string = "users"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(u *models.User) error {
	query := fmt.Sprintf("INSERT INTO %s(user_id, username,  team_name, is_active) VALUES($1, $2, $3, $4)", tableUser)
	if _, err := ur.db.Exec(query, u.UserID, u.UserName, u.TeamName, u.IsActive); err != nil {
		return err

	}
	return nil
}

func (ur *UserRepository) FindUserById(UserID string) (*models.User, error) {
	query := fmt.Sprintf("SELECT user_id, username, team_name, is_active FROM %s WHERE user_id = $1", tableUser)
	var user models.User
	err := ur.db.QueryRow(query, UserID).Scan(&user.UserID, &user.UserName, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetActiveByTeam(teamName string) ([]*models.User, error) {
	query := fmt.Sprintf("SELECT user_id, username, team_name, is_active FROM %s WHERE team_name = $1 and is_active = true ORDER BY user_id", tableUser)
	rows, err := ur.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.UserID, &u.UserName, &u.TeamName, &u.IsActive)
		if err != nil {
			messages.SendLogMessage("Cant parse user", err)
		}
		users = append(users, &u)
	}
	return users, nil
}

func (ur *UserRepository) GetAllByTeam(teamName string) ([]models.User, error) {
	query := fmt.Sprintf("SELECT user_id, username, team_name, is_active FROM %s WHERE team_name = $1 ORDER BY user_id", tableUser)
	rows, err := ur.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.UserID, &u.UserName, &u.TeamName, &u.IsActive)
		if err != nil {
			messages.SendLogMessage("Cant parse user", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (ur *UserRepository) SetUserActive(userID string, isActive bool) (*models.User, error) {
	query := fmt.Sprintf("UPDATE %s SET is_active = $2 WHERE user_id = $1 RETURNING user_id, username, team_name, is_active", tableUser)
	var user models.User
	err := ur.db.QueryRow(query, userID, isActive).Scan(&user.UserID, &user.UserName, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (ur *UserRepository) UserExist(userID string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE user_id = $1 LIMIT 1", tableUser)
	var a int

	err := ur.db.QueryRow(query, userID).Scan(&a)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ur *UserRepository) UpdateUserTeam(userID string, teamName string) error {
	query := fmt.Sprintf("UPDATE %s SET team_name = $1 WHERE user_id = $2", tableUser)
	if _, err := ur.db.Exec(query, teamName, userID); err != nil {
		return err

	}
	return nil
}
