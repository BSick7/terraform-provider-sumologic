package api

import (
	"fmt"
)

type SumologicError struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewSumologicError(res *ClientExecutorResponse) *SumologicError {
	err := &SumologicError{}
	res.BodyJSON(err)
	return err
}

func (r *SumologicError) Error() string {
	return fmt.Sprintf("[id=%s] %d %s %s", r.ID, r.Status, r.Code, r.Message)
}
