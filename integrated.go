package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
)

func TestGetFilter() {
	request, err := http.NewRequest("GET", "/filter", nil)
	if err != nil {
		logrus.Fatalf("Failed to create request: %v", err)
	}

	response := httptest.NewRecorder()

	serveFilter(response, request)

	if response.Code != http.StatusOK {
		logrus.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
	}

	contentType := response.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		logrus.Errorf("Incorrect content type. Expected: text/html; charset=utf-8, Got: %s", contentType)
	}

	logrus.Println("Integrated Test: Success")
}
