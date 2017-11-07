package api

import "net/http"

type AnonymousTransport struct {
	action func(req *http.Request) (*http.Response, error)
}

func NewAnonymousTransport(action func(req *http.Request) (*http.Response, error)) *AnonymousTransport {
	return &AnonymousTransport{action: action}
}

func (t *AnonymousTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.action(req)
}
