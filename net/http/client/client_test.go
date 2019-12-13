package client

import (
	"testing"
)

func TestUserAgent(t *testing.T) {

	ua := New()

	resp, err := ua.Request("GET", "https://www.livejournal.com", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal("wrong status code")
	}

	if string(resp.Content) == "" {
		t.Fatal("wrong respose")
	}
}
