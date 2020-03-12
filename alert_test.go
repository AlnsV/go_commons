package utils

import (
	"testing"
	"time"
)

func TestNewAlertHandler(t *testing.T) {
	t.Parallel()
	_ = NewAlertHandler(5, "https://36l9op4nsg.execute-api.us-east-2.amazonaws.com/production/")
}

func TestAlertHandler_Send(t *testing.T) {
	t.Parallel()
	alert := NewAlertHandler(3*time.Second, "https://36l9op4nsg.execute-api.us-east-2.amazonaws.com/production/")
	time.Sleep(5 * time.Second)
	err := alert.Send(
		"Test",
		"Running test on preprocessor",
		"testAlert",
	)
	if err != nil {
		t.Errorf("Failure sending email: err %s", err)
	}
	err = alert.Send(
		"Test2",
		"holhsdfsdfjk",
		"testAlert",
	)
	if err != nil {
		t.Errorf("Failure sending email: err %s", err)
	}
	alertTwo := NewAlertHandler(3*time.Second, "https://36l9op4nasdafsil872i8734hasi4ba8ys.com/asdhaskjdh/")
	time.Sleep(5 * time.Second)
	errTwo := alertTwo.Send(
		"Test",
		"Running test on preprocessor",
		"testAlert",
	)
	if errTwo == nil {
		t.Error("this host shouldn't be found")
	}
}
