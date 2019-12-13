package pixel

import (
	"net/http/httptest"
	"testing"
)

func TestSend(t *testing.T) {
	w := httptest.NewRecorder()
	Send(w)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatal("Bad status code")
	}

	if resp.Header.Get("Content-Type") != "image/gif" {
		t.Fatal("Bad content type")
	}
}
