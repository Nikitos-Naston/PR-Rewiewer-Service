package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"net/http"
)

func (H *Handler) AddTeam(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("POST team/add/", nil)

	var team models.Team
	err := json.NewDecoder(req.Body).Decode(&team)

	if err != nil {
		messages.SendLogMessage("problem with parsing JSON", err)
		messages.SendMessageJSON(writer, "INVALID_JSON")
		return
	}

	exist, err := H.TeamRepo.TeamExist(team.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB While TeamExist Operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if exist {
		messages.SendLogMessage("Operation ends by user error TEAM_EXISTS 409", nil)
		messages.SendMessageJSON(writer, "TEAM_EXISTS")
		return
	}

	err = H.TeamRepo.CreateTeam(team.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB while Create team Operaation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	for idx, memb := range team.Members {
		team.Members[idx].TeamName = team.TeamName
		memb.TeamName = team.TeamName
		exist, err = H.UserRepo.UserExist(memb.UserID)
		if err != nil {
			messages.SendLogMessage("problem with  DB while UserExist Operation", err)
			messages.SendMessageJSON(writer, "SERVER_ERROR")
			return
		}
		if exist {
			err = H.UserRepo.UpdateUserTeam(memb.UserID, memb.TeamName)
			if err != nil {
				messages.SendLogMessage("problem with  DB while UpdateUserTeam", err)
				messages.SendMessageJSON(writer, "SERVER_ERROR")
				return
			}
			_, err = H.UserRepo.SetUserActive(memb.UserID, memb.IsActive)
			if err != nil {
				messages.SendLogMessage("problem with  DB while SetIsActive", err)
				messages.SendMessageJSON(writer, "SERVER_ERROR")
				return
			}
		} else {
			err = H.UserRepo.CreateUser(&memb)
			if err != nil {
				messages.SendLogMessage("problem with  DB While CreateUser operation", err)
				messages.SendMessageJSON(writer, "SERVER_ERROR")
				return
			}
		}
	}
	messages.SendLogMessage("Team Added succesful", nil)
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(team)
}

func (H *Handler) GetTeam(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("GET team/get/", nil)

	teamName := req.URL.Query().Get("team_name")
	messages.SendLogMessage("TeamName = "+teamName, nil)
	team := models.Team{}
	team.TeamName = teamName

	exist, err := H.TeamRepo.TeamExist(team.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB while team exist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if !exist {
		messages.SendLogMessage("Operation ends by user error 404 NOT_FOUD", nil)
		messages.SendMessageJSON(writer, "NOT_FOUD")
		return
	}

	members, err := H.UserRepo.GetAllByTeam(team.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB while getting all by team operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	team.Members = members
	messages.SendLogMessage("Team Get succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(team)
}
