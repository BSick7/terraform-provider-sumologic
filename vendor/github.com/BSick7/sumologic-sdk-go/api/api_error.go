package api

import (
	"fmt"
)

type APIError struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(req *ClientExecutorRequest) *APIError {
	res := &APIError{}
	req.GetJSONBody(res)
	return res
}

func (r *APIError) Error() string {
	return fmt.Sprintf("[id=%s] %d %s %s", r.ID, r.Status, r.Code, r.Message)
}
