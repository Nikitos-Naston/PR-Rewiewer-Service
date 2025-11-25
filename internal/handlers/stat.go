package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"net/http"
)

func (H *Handler) ShowStat(writer http.ResponseWriter, req *http.Request) {
	messages.SendLogMessage("Getting Stat start", nil)
	globalStats, err := H.StatRepo.GetPRStats()
	if err != nil {
		messages.SendLogMessage("problem with  DB while GettingALLStats operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	UserStats, err := H.StatRepo.GetUserStats()
	if err != nil {
		messages.SendLogMessage("problem with  DB while GettingUserStats operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	response := models.StatsResponnse{
		GlobalStat: *globalStats,
		UserStat:   UserStats,
	}
	messages.SendLogMessage("PR added succesful", nil)
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(response)
}
