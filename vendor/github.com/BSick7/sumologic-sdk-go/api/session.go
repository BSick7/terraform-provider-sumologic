package api

/// All implementations of Session are expected to use the following environment variables:
///   SUMO_ACCESS_ID
///   SUMO_ACCESS_KEY
///   SUMO_DEBUG (turns on request/response logging)

import (
	"net/http"
	"net/url"
)

type Session interface {
	SetAddress(address string)
	SetCredentials(accessID, accessKey string)
	EndpointURL(endpoint string) *url.URL
	CreateTransport() http.RoundTripper
}
