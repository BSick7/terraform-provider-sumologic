package api

import (
	"net/http"
	"testing"
)

const (
	discoverResponse = `
{
  "status" : 301,
  "id" : "B5DT0-EJ0IP-85W4L",
  "code" : "moved",
  "message" : "The requested resource SHOULD be accessed through returned URI in Location Header."
}`
)

func TestSession_Discover(t *testing.T) {
	session := NewMockSession(false)
	session.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Location", "http://redirected.api.sumologic.com")
			http.Error(w, discoverResponse, http.StatusMovedPermanently)
		} else {
			http.Error(w, "unexpected request", http.StatusInternalServerError)
		}
	}))

	session.Discover()

	got := session.EndpointURL("/path").String()
	want := "http://redirected.api.sumologic.com/path"
	if got != want {
		t.Errorf("mismatched endpoints, got %s, want %s", got, want)
	}
}
