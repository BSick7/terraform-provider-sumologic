package api

import (
	"net/http"
)

type ClientExecutor struct {
	session Session
	client  *http.Client
}

func NewClientExecutor(session Session, client *http.Client) *ClientExecutor {
	return &ClientExecutor{
		session: session,
		client:  client,
	}
}

func (c *ClientExecutor) NewRequest() (*ClientExecutorRequest, error) {
	return NewClientExecutorRequest(c.session, c.client)
}
