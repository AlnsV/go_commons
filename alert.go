package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Alert interface {
	Send(subject, message, alertTopic string) error
}

// AlertHandler manages the timeout for generating an alert
type AlertHandler struct {
	activationTimes int
	requestUrl      string
	lastSubject     string
	lastMessage     string
	lastAlert       time.Time
	timeout         time.Duration
}

func NewAlertHandler(timeout time.Duration, url string) *AlertHandler {
	return &AlertHandler{
		requestUrl:      url,
		activationTimes: 0,
		timeout:         timeout,
	}
}

// AlertHandler_Send this method helps into sending an alert email
func (a *AlertHandler) Send(subject, message, alertTopic string) error {
	// Verifies if the alert activated too recently
	if a.activationTimes != 0 && time.Now().UTC().Sub(a.lastAlert) < a.timeout {
		a.activationTimes += 1
		return nil
	}

	a.activationTimes += 1
	a.lastAlert = time.Now().UTC()
	a.lastMessage = message
	a.lastSubject = subject

	// Formatting the request body with the relevant info for our request
	reqBody, _ := json.Marshal(map[string]string{
		"subject": subject,
		"mailBody": fmt.Sprintf(
			"%s\nActivated: %d times.\nTimestamp:%s",
			message,
			a.activationTimes,
			a.lastAlert.Format(time.RFC1123),
		),
		"topic": fmt.Sprintf("arn:aws:sns:us-east-2:406163507813:%s", alertTopic),
	})

	// Execution of the post requests with the follow URL
	_, err := http.Post(a.requestUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failure sending Post request for sending email alert, err: %s", err)
	}
	return nil
}
