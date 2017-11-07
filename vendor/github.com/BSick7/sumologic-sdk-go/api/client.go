package api

import (
	"net/http"
	"strings"
)

type Client struct {
	session  Session
	executor *ClientExecutor
}

func NewClient(session Session) *Client {
	return &Client{
		session:  session,
		executor: NewClientExecutor(session, createHttpClient(session)),
	}
}

func createHttpClient(session Session) *http.Client {
	return &http.Client{
		Transport: session.CreateTransport(),
	}
}

func (c *Client) Discover() {
	req, _ := c.executor.NewRequest()
	req.SetEndpoint("/")
	req.Put()
	location := req.GetResponseHeader("Location")
	if location != "" {
		c.session.SetAddress(strings.TrimRight(location, "/"))
	}
}

func (c *Client) Collectors() *Collectors {
	return NewCollectors(c.executor)
}
