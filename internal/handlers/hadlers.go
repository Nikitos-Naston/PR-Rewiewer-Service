package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/repository"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	TeamRepo *repository.TeamRepository
	UserRepo *repository.UserRepository
	PRRepo   *repository.PRRepository
}

func NewHandler(TeamRepo *repository.TeamRepository, UserRepo *repository.UserRepository, PRRepo *repository.PRRepository) *Handler {
	return &Handler{
		TeamRepo: TeamRepo,
		UserRepo: UserRepo,
		PRRepo:   PRRepo,
	}
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (H *Handler) Status(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "Check api helth", nil)
	messages.SendMessageJSON(writer, 200, "Health")
}
