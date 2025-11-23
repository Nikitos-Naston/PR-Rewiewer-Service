package messages

import "github.com/sirupsen/logrus"

var logger = *logrus.New()

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func SendLogMessage(level logrus.Level, message string, err error) {
	logger.SetLevel(level)
	if err != nil {
		logger.Println(message, err)
	} else {
		logger.Println(message)
	}
}
