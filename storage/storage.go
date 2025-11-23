package storage

import (
	"PRreviewService/config"
	"PRreviewService/internal/messages"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

func New(cfg config.Config) (*sql.DB, error) {
	InfoDB := fmt.Sprintf("host = %s port = %s user = %s password = %s dbname = %s sslmode = disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	messages.SendLogMessage(logrus.InfoLevel, InfoDB, nil)

	db, err := sql.Open("postgres", InfoDB)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
