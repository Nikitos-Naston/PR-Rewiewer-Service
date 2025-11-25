package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/repository"
	"encoding/json"
	"net/http"
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
	messages.SendLogMessage("Check api helth", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode("Healthy")
}
