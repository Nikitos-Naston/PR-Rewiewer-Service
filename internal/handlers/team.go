package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (H *Handler) AddTeam(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "POST team/add/", nil)

	var team models.Team
	err := json.NewDecoder(req.Body).Decode(&team)

	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with parsing JSON", err)
		messages.SendMessageJSON(writer, 400, "Problem with your JSON file")
		return
	}

	exist, err := H.TeamRepo.TeamExist(team.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB While TeamExist Operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 400, "TEAM_EXISTS")
		return
	}

	err = H.TeamRepo.CreateTeam(team.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while Create team Operaation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
	}

	for idx, memb := range team.Members {
		team.Members[idx].TeamName = team.TeamName
		memb.TeamName = team.TeamName
		exist, err = H.UserRepo.UserExist(memb.UserID)
		if err != nil {
			messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while UserExist Operation", err)
			messages.SendMessageJSON(writer, 501, "Server temporary broke")
			return
		}
		if exist {
			err = H.UserRepo.UpdateUserTeam(memb.UserID, memb.TeamName)
			if err != nil {
				messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while UpdateUserTeam", err)
				messages.SendMessageJSON(writer, 501, "Server temporary broke")
				return
			}
		} else {
			err = H.UserRepo.CreateUser(&memb)
			if err != nil {
				messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB While CreateUser operation", err)
				messages.SendMessageJSON(writer, 501, "Server temporary broke")
				return
			}
		}
	}
	messages.SendLogMessage(logrus.InfoLevel, "Team Added succesful", nil)
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(team)
}

func (H *Handler) GetTeam(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "GET team/get/", nil)

	teamName := req.PathValue("team_name")
	messages.SendLogMessage(logrus.DebugLevel, "TeamName = "+teamName, nil)
	team := models.Team{}
	team.TeamName = teamName

	exist, err := H.TeamRepo.TeamExist(team.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while team exist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "NOT_FOUD")
		return
	}

	members, err := H.UserRepo.GetAllByTeam(team.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while getting all by team operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	team.Members = members
	messages.SendLogMessage(logrus.InfoLevel, "Team Get succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(team)
}
