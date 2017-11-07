package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type ClientExecutorRequest struct {
	debug   bool
	session Session
	client  *http.Client
	req     *http.Request
	reqBody *bytes.Buffer
	res     *http.Response
}

func NewClientExecutorRequest(session Session, client *http.Client) (*ClientExecutorRequest, error) {
	buf := bytes.NewBuffer(nil)
	req, err := http.NewRequest("GET", "/", buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	return &ClientExecutorRequest{
		debug:   os.Getenv("SUMO_DEBUG") == "1",
		session: session,
		client:  client,
		req:     req,
		reqBody: buf,
	}, nil
}

func (r *ClientExecutorRequest) SetEndpoint(endpoint string) *ClientExecutorRequest {
	r.req.URL = r.session.EndpointURL(endpoint)
	return r
}

func (r *ClientExecutorRequest) SetQuery(params url.Values) *ClientExecutorRequest {
	r.req.URL.RawQuery = params.Encode()
	return r
}

func (r *ClientExecutorRequest) SetRequestHeader(key string, value string) *ClientExecutorRequest {
	r.req.Header.Set(key, value)
	return r
}

func (r *ClientExecutorRequest) SetJSONBody(input interface{}) error {
	encoder := json.NewEncoder(r.reqBody)
	if err := encoder.Encode(input); err != nil {
		return fmt.Errorf("error encoding json body: %s", err)
	}
	r.req.Header.Set("Content-Type", "application/json")
	return nil
}

func (r *ClientExecutorRequest) Get() error {
	r.req.Method = "GET"
	return r.do()
}

func (r *ClientExecutorRequest) Post() error {
	r.req.Method = "POST"
	return r.do()
}

func (r *ClientExecutorRequest) Put() error {
	r.req.Method = "PUT"
	return r.do()
}

func (r *ClientExecutorRequest) Delete() error {
	r.req.Method = "DELETE"
	return r.do()
}

func (r *ClientExecutorRequest) do() error {
	if r.debug {
		raw, _ := httputil.DumpRequestOut(r.req, true)
		log.Println(string(raw))
	}

	res, err := r.client.Do(r.req)
	r.res = res
	if r.debug && res != nil {
		raw, _ := httputil.DumpResponse(res, false)
		log.Println(string(raw))
	}
	if err != nil {
		return fmt.Errorf("error requesting %s %s: %s", r.req.Method, r.req.URL, err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return NewAPIError(r)
	}

	return nil
}

func (r *ClientExecutorRequest) GetJSONBody(out interface{}) error {
	decoder := json.NewDecoder(r.res.Body)
	defer r.res.Body.Close()
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("error decoding json body: %s", err)
	}
	return nil
}

func (r *ClientExecutorRequest) GetStringBody() (string, error) {
	defer r.res.Body.Close()
	raw, err := ioutil.ReadAll(r.res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %s", err)
	}
	return string(raw), nil
}

func (r *ClientExecutorRequest) GetRawBody() ([]byte, error) {
	defer r.res.Body.Close()
	raw, err := ioutil.ReadAll(r.res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %s", err)
	}
	return raw, nil
}

func (r *ClientExecutorRequest) GetResponseHeader(key string) string {
	return r.res.Header.Get(key)
}
