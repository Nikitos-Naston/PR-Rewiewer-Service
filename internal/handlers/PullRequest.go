package handlers

import (
	"PRreviewService/internal"
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/sirupsen/logrus"
)

func (H *Handler) CreatePR(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "POST /pullRequest/create", nil)

	var PR *models.PR
	err := json.NewDecoder(req.Body).Decode(&PR)

	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with parsing JSON", err)
		messages.SendMessageJSON(writer, 400, "Problem with your JSON file")
		return
	}

	exist, err := H.UserRepo.UserExist(PR.AuthorID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "NOT_FOUD")
		return
	}

	user, err := H.UserRepo.FindUserById(PR.AuthorID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while FindUserbyID operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	exist, err = H.TeamRepo.TeamExist(user.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while TEamExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "NOT_FOUD")
		return
	}

	exist, err = H.PRRepo.PRExist(PR.ID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while PRExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	fmt.Println(exist, PR, PR.Name, PR.ID)
	if exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 409, "PR_EXISTS")
		return
	}

	users, err := H.UserRepo.GetActiveByTeam(user.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while GetActiveByTeam operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	PR, err = H.PRRepo.CreatePR(*PR)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while CreateRP operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	cnt := 0
	for _, u := range users {
		if u.UserID != user.UserID {
			PR.AssignedReviewers = append(PR.AssignedReviewers, u.UserID)
			cnt += 1
			err = H.PRRepo.CreatePRRewie(PR.ID, u.UserID)
			if err != nil {
				messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while CreatePRRewie operation", err)
				messages.SendMessageJSON(writer, 501, "Server temporary broke")
				return
			}
		}
		if cnt == 2 {
			break
		}
	}
	AssignRewiers := PR.AssignedReviewers
	PR, err = H.PRRepo.GetPRByID(PR.ID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while GetPRby operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	PR.AssignedReviewers = AssignRewiers
	messages.SendLogMessage(logrus.InfoLevel, "PR added succesful", nil)
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(PR)
}

func (H *Handler) MergePR(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "POST /pullRequest/merge", nil)

	var PR *models.PR
	err := json.NewDecoder(req.Body).Decode(&PR)

	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with parsing JSON", err)
		messages.SendMessageJSON(writer, 400, "Problem with your JSON file")
		return
	}

	exist, err := H.PRRepo.PRExist(PR.ID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while PRExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "NOT_FOUND")
		return
	}

	err = H.PRRepo.MERGEPR(PR.ID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while Merge operations", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	PR, err = H.PRRepo.GetPRByID(PR.ID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while GetPRbyID operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	for _, id := range PR.AssignedReviewers {
		err = H.PRRepo.CreatePRRewie(PR.ID, id)
		if err != nil {
			messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while CreatePRRewie operation", err)
			messages.SendMessageJSON(writer, 501, "Server temporary broke")
			return
		}
	}

	messages.SendLogMessage(logrus.InfoLevel, "PR merge succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(PR)
}

func (H *Handler) RessignRewiewer(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage(logrus.InfoLevel, "POST /pullRequest/ressign", nil)

	var RRpr models.ReassignRequest
	err := json.NewDecoder(req.Body).Decode(&RRpr)

	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with parsing JSON", err)
		messages.SendMessageJSON(writer, 400, "Problem with your JSON file")
		return
	}

	exist, err := H.UserRepo.UserExist(RRpr.OldUserID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "NOT_FOUD")
		return
	}

	exist, err = H.PRRepo.PRExist(RRpr.PullRequestID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while PRExist operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	if !exist {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 404, "NOT_FOUND")
		return
	}

	PR, err := H.PRRepo.GetPRByID(RRpr.PullRequestID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while GetPRbyID operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	PR.AssignedReviewers, err = H.PRRepo.GetReviewers(PR.ID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while GetReviewers operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}
	if PR.Status == "MERGED" {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 409, "PR_MERGED")
		return
	}

	found := slices.Contains(PR.AssignedReviewers, RRpr.OldUserID)

	if found {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 409, "NOT_ASSIGNED")
		return
	}

	user, err := H.UserRepo.FindUserById(PR.AuthorID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while FindUserbyID operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	users, err := H.UserRepo.GetActiveByTeam(user.TeamName)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while GetActiveByTeam operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	err = H.PRRepo.DeletePR(RRpr.PullRequestID, RRpr.OldUserID)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while DeletePR operation", err)
		messages.SendMessageJSON(writer, 501, "Server temporary broke")
		return
	}

	chosen := false
	chosenID := ""
	for _, u := range users {
		if u.UserID != PR.AuthorID && !(slices.Contains(PR.AssignedReviewers, u.UserID)) {
			PR.AssignedReviewers = append(PR.AssignedReviewers, u.UserID)
			chosenID = u.UserID
			chosen = true
			err = H.PRRepo.CreatePRRewie(PR.ID, u.UserID)
			if err != nil {
				messages.SendLogMessage(logrus.ErrorLevel, "problem with  DB while CreatePRRewie operation", err)
				messages.SendMessageJSON(writer, 501, "Server temporary broke")
				return
			}
			break
		}
	}
	if !chosen {
		messages.SendLogMessage(logrus.InfoLevel, "Operation ends by user error", nil)
		messages.SendMessageJSON(writer, 409, "NO_CANDIDATE")
		return
	}

	PR.AssignedReviewers = internal.DeleteElemInSlice(PR.AssignedReviewers, RRpr.OldUserID)

	result := models.AnswerbyReassignRequest{
		PR:    PR,
		NewId: chosenID,
	}

	messages.SendLogMessage(logrus.InfoLevel, "PR Ressign succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(result)
}
