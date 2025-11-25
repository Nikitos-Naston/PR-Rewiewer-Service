package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"net/http"
)

func (H *Handler) SetUserActive(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("POST /users/setIsActive", nil)

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		messages.SendLogMessage("problem with parsing JSON", err)
		messages.SendMessageJSON(writer, "INVALID_JSON")
		return
	}

	exist, err := H.UserRepo.UserExist(user.UserID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if !exist {
		messages.SendLogMessage("Operation ends by user error USER_NOT_FOUD 404", nil)
		messages.SendMessageJSON(writer, "USER_NOT_FOUD")
		return
	}

	u, err := H.UserRepo.SetUserActive(user.UserID, user.IsActive)
	if err != nil {
		messages.SendLogMessage("problem with  DB while SetUserActive Operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	messages.SendLogMessage("User set activity succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(u)
}

func (H *Handler) GetRewiesByUser(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("GET /users/getReview", nil)

	id := req.URL.Query().Get("user_id")

	exist, err := H.UserRepo.UserExist(id)
	if err != nil {
		messages.SendLogMessage("problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if !exist {
		messages.SendLogMessage("Operation ends by user error USER_NOT_FOUD 404", nil)
		messages.SendMessageJSON(writer, "USER_NOT_FOUD")
		return
	}

	prSliceID, err := H.PRRepo.GetRRByUserID(id)
	if err != nil {
		messages.SendLogMessage("problem with  DB while getuserbyID  operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	prSlice := make([]*models.PRsmall, 0)
	for _, id := range prSliceID {
		PR, err := H.PRRepo.GetPRByID(id)
		if err != nil {
			messages.SendLogMessage("problem with  DB while getuserbyID  operation", err)
			messages.SendMessageJSON(writer, "SERVER_ERROR")
			return
		}
		prSlice = append(prSlice, &models.PRsmall{
			ID:       PR.ID,
			Name:     PR.Name,
			AuthorID: PR.AuthorID,
			Status:   PR.Status,
		})
	}

	messages.SendLogMessage("Get rewiews by ID succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(prSlice)
}
