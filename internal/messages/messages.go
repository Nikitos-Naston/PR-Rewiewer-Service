package messages

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var logger = *logrus.New()

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// TODO поменять везде LOGGER
func SendLogMessage(level logrus.Level, message string, err error) {
	logger.SetLevel(level)
	if err != nil {
		logger.Error(message, err)
	} else {
		logger.Info(message)
	}
}

func SendMessageJSON(w http.ResponseWriter, statusCode int, message string) {
	msg := Message{
		StatusCode: statusCode,
		Message:    message,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(msg)
}
