package api

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCollectors_Create(t *testing.T) {
	session := NewMockSession(true)
	client := NewClient(session)

	response := `{
  "collector": {
    "id": 100772723,
    "name": "My Hosted Collector",
    "description": "An example Hosted Collector",
    "category": "HTTP Collection",
    "timeZone": "UTC",
    "links": [
      {
        "rel": "sources",
        "href": "/v1/collectors/100772723/sources"
      }
    ],
    "collectorType": "Hosted",
    "collectorVersion": "",
    "lastSeenAlive": 1476818195411,
    "alive": true
  }
}`

	session.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/v1/collectors" && r.Header.Get("Content-Type") == "application/json" {
			w.Write([]byte(response))
		} else {
			http.Error(w, "unexpected request", http.StatusInternalServerError)
		}
	}))

	got, err := client.Collectors().Create(&CollectorCreate{
		CollectorType: "Hosted",
		Name:          "My Hosted Collector",
		Description:   "An example Hosted Collector",
		Category:      "HTTP Collection",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := &Collector{
		ID:            100772723,
		CollectorType: "Hosted",
		Name:          "My Hosted Collector",
		Description:   "An example Hosted Collector",
		Category:      "HTTP Collection",
		Links: []CollectorLink{
			{
				Rel:  "sources",
				Href: "/v1/collectors/100772723/sources",
			},
		},
		TimeZone:         "UTC",
		CollectorVersion: "",
		LastSeenAlive:    1476818195411,
		Alive:            true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mismatched result, got\n%+v\nwant\n%+v", got, want)
	}
}

func TestCollectors_Update(t *testing.T) {
	session := NewMockSession(true)
	client := NewClient(session)

	response := `{
  "collector": {
    "id": 100772723,
    "name": "My Hosted Collector",
    "description": "An example Hosted Collector",
    "category": "HTTP Collection",
    "timeZone": "UTC",
    "links": [
      {
        "rel": "sources",
        "href": "/v1/collectors/100772723/sources"
      }
    ],
    "collectorType": "Hosted",
    "collectorVersion": "",
    "lastSeenAlive": 1476818195411,
    "alive": true
  }
}`

	etag := "stub-etag"
	session.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/v1/collectors/100772723" {
			w.Header().Set("ETag", etag)
			w.Write([]byte(response))
		} else if r.Method == "PUT" && r.URL.Path == "/api/v1/collectors/100772723" && r.Header.Get("Content-Type") == "application/json" {
			if r.Header.Get("If-Match") != etag {
				t.Errorf("mismatched If-Match, got %s, want %s", r.Header.Get("If-Match"), etag)
			}
			w.Write([]byte(response))
		} else {
			http.Error(w, "unexpected request", http.StatusInternalServerError)
		}
	}))

	got, err := client.Collectors().Update(&Collector{
		ID:            100772723,
		CollectorType: "Hosted",
		Name:          "My Hosted Collector",
		Description:   "An example Hosted Collector",
		Category:      "HTTP Collection",
		Links: []CollectorLink{
			{
				Rel:  "sources",
				Href: "/v1/collectors/100772723/sources",
			},
		},
		TimeZone:         "UTC",
		CollectorVersion: "",
		LastSeenAlive:    1476818195411,
		Alive:            true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := &Collector{
		ID:            100772723,
		CollectorType: "Hosted",
		Name:          "My Hosted Collector",
		Description:   "An example Hosted Collector",
		Category:      "HTTP Collection",
		Links: []CollectorLink{
			{
				Rel:  "sources",
				Href: "/v1/collectors/100772723/sources",
			},
		},
		TimeZone:         "UTC",
		CollectorVersion: "",
		LastSeenAlive:    1476818195411,
		Alive:            true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mismatched result, got\n%+v\nwant\n%+v", got, want)
	}
}
