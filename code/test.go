package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnitPostLogin(t *testing.T) {
	payload := `{"name":"testuser","password":"password"}`
	req, err := http.NewRequest("POST", "/postLogin", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(postLogin)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")

	assert.Contains(t, rr.Body.String(), "success", "Response should contain success message")
}

func TestIntegrationAddVideoByUser(t *testing.T) {
	payload := `{"title":"Test Title","author":"Test Author","likes":"10","comments":"Nice video","date":"2024-03-03","imagepath":"test.jpg"}`
	req, err := http.NewRequest("POST", "/addVideoByUser", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addVideoByUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code, "Status code should be StatusFound")
}

func TestE2EPostFilter(t *testing.T) {
	req, err := http.NewRequest("POST", "/postFilter?page=1&pageSize=10&sortOrder=1", strings.NewReader(`{"title":"Test","author":"Author"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postFilter)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
	assert.Contains(t, rr.Body.String(), "Test Title", "Response should contain test title")
}
