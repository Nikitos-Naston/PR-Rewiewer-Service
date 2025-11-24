package handlers

import (
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (H *Handler) SetUserActive(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "POST /users/setIsActive", nil)

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with parsing JSON", err)
		messages.SendMessageJSON(writer, 400, "Problem with your JSON file")
		return
	}

	exist, err := H.UserRepo.UserExist(user.UserID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "USER_NOT_FOUD")
		return
	}

	u, err := H.UserRepo.SetUserActive(user.UserID, user.IsActive)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while SetUserActive Operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	messages.SendLogMessage(logrus.InfoLevel, "User set activity succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(u)
}

func (H *Handler) GetRewiesByUser(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "GET /users/getReview", nil)

	id := req.PathValue("user_id")

	exist, err := H.UserRepo.UserExist(id)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "USER_NOT_FOUD")
		return
	}

	prSliceID, err := H.PRRepo.GetRRByUserID(id)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while getuserbyID  operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	prSlice := make([]*models.PR, 0)
	for _, id := range prSliceID {
		PR, err := H.PRRepo.GetPRByID(id)
		if err != nil {
			messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while getuserbyID  operation", err)
			messages.SendMessageJSON(writer, 501, "Server temporary broke")
			return
		}
		prSlice = append(prSlice, PR)
	}
	// TODO converting PR to PRSMALL
	messages.SendLogMessage(logrus.InfoLevel, "Get rewiews by ID succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(prSlice)
}
