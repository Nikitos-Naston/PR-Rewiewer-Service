package handlers

import (
	"PRreviewService/internal"
	"PRreviewService/internal/messages"
	"PRreviewService/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

func (H *Handler) CreatePR(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("POST /pullRequest/create", nil)

	var PR *models.PR
	err := json.NewDecoder(req.Body).Decode(&PR)

	if err != nil {
		messages.SendLogMessage("problem with parsing JSON", err)
		messages.SendMessageJSON(writer, "INVALID_JSON")
		return
	}

	exist, err := H.UserRepo.UserExist(PR.AuthorID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if !exist {
		messages.SendLogMessage("Operation ends by user error NOT_FOUD 404 User", nil)
		messages.SendMessageJSON(writer, "NOT_FOUD")
		return
	}

	user, err := H.UserRepo.FindUserById(PR.AuthorID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while FindUserbyID operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	if !user.IsActive {
		messages.SendLogMessage("Operation ends by user error USER_MUST_BE_ACTIVE 409 User", nil)
		messages.SendMessageJSON(writer, "AUTHOR_INACTIVE")
		return
	}
	exist, err = H.TeamRepo.TeamExist(user.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB while TEamExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	if !exist {
		messages.SendLogMessage("Operation ends by user error NOT_FOUD 404 Team", nil)
		messages.SendMessageJSON(writer, "NOT_FOUD")
		return
	}

	exist, err = H.PRRepo.PRExist(PR.ID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while PRExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	fmt.Println(exist, PR, PR.Name, PR.ID)
	if exist {
		messages.SendLogMessage("Operation ends by user error PR_EXISTS 409", nil)
		messages.SendMessageJSON(writer, "PR_EXISTS")
		return
	}

	users, err := H.UserRepo.GetActiveByTeam(user.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB while GetActiveByTeam operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	PR, err = H.PRRepo.CreatePR(*PR)
	if err != nil {
		messages.SendLogMessage("problem with  DB while CreateRP operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	cnt := 0
	for _, u := range users {
		if u.UserID != user.UserID {
			PR.AssignedReviewers = append(PR.AssignedReviewers, u.UserID)
			cnt += 1
			err = H.PRRepo.CreatePRRewie(PR.ID, u.UserID)
			if err != nil {
				messages.SendLogMessage("problem with  DB while CreatePRRewie operation", err)
				messages.SendMessageJSON(writer, "SERVER_ERROR")
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
		messages.SendLogMessage("problem with  DB while GetPRby operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	PR.AssignedReviewers = AssignRewiers
	messages.SendLogMessage("PR added succesful", nil)
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(PR)
}

func (H *Handler) MergePR(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("POST /pullRequest/merge", nil)

	var PR *models.PR
	err := json.NewDecoder(req.Body).Decode(&PR)

	if err != nil {
		messages.SendLogMessage("problem with parsing JSON", err)
		messages.SendMessageJSON(writer, "INVALID_JSON")
		return
	}

	exist, err := H.PRRepo.PRExist(PR.ID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while PRExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	if !exist {
		messages.SendLogMessage("Operation ends by user error NOT_FOUND 404 PR", nil)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if PR.Status == "OPEN" {
		err = H.PRRepo.MERGEPR(PR.ID)
		if err != nil {
			messages.SendLogMessage("problem with  DB while Merge operations", err)
			messages.SendMessageJSON(writer, "SERVER_ERROR")
			return
		}
	}
	PR, err = H.PRRepo.GetPRByID(PR.ID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while GetPRbyID operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	for _, id := range PR.AssignedReviewers {
		err = H.PRRepo.CreatePRRewie(PR.ID, id)
		if err != nil {
			messages.SendLogMessage("problem with  DB while CreatePRRewie operation", err)
			messages.SendMessageJSON(writer, "SERVER_ERROR")
			return
		}
	}
	PR.Status = "MERGED"
	messages.SendLogMessage("PR merge succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(PR)
}

func (H *Handler) RessignRewiewer(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	messages.SendLogMessage("POST /pullRequest/ressign", nil)

	var RRpr models.ReassignRequest
	err := json.NewDecoder(req.Body).Decode(&RRpr)

	if err != nil {
		messages.SendLogMessage("problem with parsing JSON", err)
		messages.SendMessageJSON(writer, "INVALID_JSON")
		return
	}

	exist, err := H.UserRepo.UserExist(RRpr.OldUserID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while UserExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if !exist {
		messages.SendLogMessage("Operation ends by user error NOT_FOUD 404 USER", nil)
		messages.SendMessageJSON(writer, "NOT_FOUD")
		return
	}

	exist, err = H.PRRepo.PRExist(RRpr.PullRequestID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while PRExist operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	if !exist {
		messages.SendLogMessage("Operation ends by user error NOT_FOUND 404 PR", nil)
		messages.SendMessageJSON(writer, "NOT_FOUND")
		return
	}

	PR, err := H.PRRepo.GetPRByID(RRpr.PullRequestID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while GetPRbyID operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	PR.AssignedReviewers, err = H.PRRepo.GetReviewers(PR.ID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while GetReviewers operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}
	if PR.Status == "MERGED" {
		messages.SendLogMessage("Operation ends by user error PR_MERGED 409", nil)
		messages.SendMessageJSON(writer, "PR_MERGED")
		return
	}

	found := slices.Contains(PR.AssignedReviewers, RRpr.OldUserID)

	if !found {
		messages.SendLogMessage("Operation ends by user error NOT_ASSIGNED 409", nil)
		messages.SendMessageJSON(writer, "NOT_ASSIGNED")
		return
	}

	user, err := H.UserRepo.FindUserById(PR.AuthorID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while FindUserbyID operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	users, err := H.UserRepo.GetActiveByTeam(user.TeamName)
	if err != nil {
		messages.SendLogMessage("problem with  DB while GetActiveByTeam operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
		return
	}

	err = H.PRRepo.DeletePR(RRpr.PullRequestID, RRpr.OldUserID)
	if err != nil {
		messages.SendLogMessage("problem with  DB while DeletePR operation", err)
		messages.SendMessageJSON(writer, "SERVER_ERROR")
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
				messages.SendLogMessage("problem with  DB while CreatePRRewie operation", err)
				messages.SendMessageJSON(writer, "SERVER_ERROR")
				return
			}
			break
		}
	}
	if !chosen {
		messages.SendLogMessage("Operation ends by user error NO_CANDIDATE 409", nil)
		messages.SendMessageJSON(writer, "NO_CANDIDATE")
		return
	}

	PR.AssignedReviewers = internal.DeleteElemInSlice(PR.AssignedReviewers, RRpr.OldUserID)

	result := models.AnswerbyReassignRequest{
		PR:    PR,
		NewId: chosenID,
	}

	messages.SendLogMessage("PR Ressign succesful", nil)
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(result)
}
