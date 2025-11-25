package messages

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var logger = *logrus.New()

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type ErrorResponse struct {
	Error AppError `json:"error"`
}

func SendLogMessage(message string, err error) {
	if err != nil {
		logger.Error(message, err)
	} else {
		logger.Info(message)
	}
}

func SendMessageJSON(w http.ResponseWriter, statusCode string) {

	Response := ErrorResponse{
		Error: AppError{
			Code:    statusCode,
			Message: GetErrorMessage(statusCode),
		},
	}

	w.WriteHeader(GetHTTPStatusCode(statusCode))
	json.NewEncoder(w).Encode(Response)
}

func GetHTTPStatusCode(errorCode string) int {
	switch errorCode {

	case "INVALID_JSON":
		return 400

	case "NOT_FOUND":
		return 404

	case "TEAM_EXISTS":
		return 409
	case "PR_EXISTS":
		return 409
	case "PR_MERGED":
		return 409
	case "NOT_ASSIGNED":
		return 409
	case "NO_CANDIDATE":
		return 409
	case "AUTHOR_INACTIVE":
		return 409

	case "SERVER_ERROR":
		return 500

	default:
		return 500
	}
}

func GetErrorMessage(errorCode string) string {
	switch errorCode {
	case "INVALID_JSON":
		return "Invalid JSON file"
	case "NOT_FOUND":
		return "Resource not found"
	case "TEAM_EXISTS":
		return "Team already exists"
	case "PR_EXISTS":
		return "Pull request already exists"
	case "PR_MERGED":
		return "Pull request is already merged"
	case "NOT_ASSIGNED":
		return "User is not assigned as reviewer"
	case "NO_CANDIDATE":
		return "No available reviewers to assign"
	case "AUTHOR_INACTIVE":
		return "Author must be active to create PR"
	case "SERVER_ERROR":
		return "Service is temporarily unavailable"
	default:
		return "Unknown error"
	}
}
