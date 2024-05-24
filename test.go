package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"testing"
)

func TestServerEndpoints(t *testing.T) {
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{"--headless"},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer wd.Quit()

	// Navigate to server endpoints and test responses
	endpoints := []string{"/", "/index", "/addVideo", "/notfound", "/login", "/register", "/createUser", "/postLogin", "/registered", "/forgot", "/filter", "/postFilter", "/starter", "/addVideoByUser"}
	for _, endpoint := range endpoints {
		err := wd.Get("http://localhost:8888" + endpoint)
		if err != nil {
			t.Fatalf("Failed to load page %s: %v", endpoint, err)
		}
		title, err := wd.Title()
		if err != nil {
			t.Fatalf("Failed to get page title for %s: %v", endpoint, err)
		}
		assert.NotEmpty(t, title, "Title should not be empty for %s", endpoint)
	}
}
