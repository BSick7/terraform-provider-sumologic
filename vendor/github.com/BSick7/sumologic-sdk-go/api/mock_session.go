package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type MockSession struct {
	impl    *SessionImpl
	server  *httptest.Server
	handler http.HandlerFunc
}

func NewMockSession(checkAuth bool) *MockSession {
	s := &MockSession{
		handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkAuth {
			if user, pw, ok := r.BasicAuth(); !ok {
				http.Error(w, "Full authentication is required to access this resource", http.StatusUnauthorized)
				return
			} else if user != s.impl.accessID || pw != s.impl.accessKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		s.handler(w, r)
	}))

	s.impl = &SessionImpl{
		accessID:  "mockaccessid",
		accessKey: "mockaccesskey",
		address:   fmt.Sprint(s.server.URL, "/api/v1"),
	}

	return s
}

func (s *MockSession) Handle(handler http.HandlerFunc) {
	s.handler = handler
}

func (s *MockSession) Discover() {
	s.impl.Discover()
}

func (s *MockSession) SetAddress(address string) {
	s.impl.SetAddress(address)
}

func (s *MockSession) SetCredentials(accessID, accessKey string) {
	s.impl.SetCredentials(accessID, accessKey)
}

func (s *MockSession) Address() string {
	return s.impl.Address()
}

func (s *MockSession) EndpointURL(endpoint string) *url.URL {
	return s.impl.EndpointURL(endpoint)
}

func (s *MockSession) CreateTransport() http.RoundTripper {
	return s.impl.CreateTransport()
}
