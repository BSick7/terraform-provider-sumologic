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

func IsObjectFound(_ interface{}, err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if serr, ok := err.(*SumologicError); ok && serr.Status == 404 {
		return false, nil
	}
	return false, err
}
