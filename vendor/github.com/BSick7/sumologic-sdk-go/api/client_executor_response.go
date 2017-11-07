package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ClientExecutorResponse struct {
	res *http.Response
}

func (r *ClientExecutorResponse) StatusCode() int {
	return r.res.StatusCode
}

func (r *ClientExecutorResponse) BodyJSON(out interface{}) error {
	decoder := json.NewDecoder(r.res.Body)
	defer r.res.Body.Close()
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("error decoding json body: %s", err)
	}
	return nil
}

func (r *ClientExecutorResponse) BodyString() (string, error) {
	defer r.res.Body.Close()
	raw, err := ioutil.ReadAll(r.res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %s", err)
	}
	return string(raw), nil
}

func (r *ClientExecutorResponse) BodyRaw() ([]byte, error) {
	defer r.res.Body.Close()
	raw, err := ioutil.ReadAll(r.res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %s", err)
	}
	return raw, nil
}

func (r *ClientExecutorResponse) Header(key string) string {
	return r.res.Header.Get(key)
}
