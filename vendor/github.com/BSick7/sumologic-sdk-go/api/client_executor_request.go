package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type ClientExecutorRequest struct {
	debug    bool
	session  Session
	client   *http.Client
	endpoint *url.URL
	body     []byte
	headers  map[string]string
}

func NewClientExecutorRequest(session Session, client *http.Client) (*ClientExecutorRequest, error) {
	return &ClientExecutorRequest{
		debug:    os.Getenv("SUMO_DEBUG") == "1",
		session:  session,
		client:   client,
		endpoint: session.EndpointURL("/"),
		body:     nil,
		headers:  map[string]string{},
	}, nil
}

func (r *ClientExecutorRequest) SetEndpoint(endpoint string) *ClientExecutorRequest {
	query := ""
	if r.endpoint != nil {
		query = r.endpoint.RawQuery
	}
	eurl := r.session.EndpointURL(endpoint)
	if eurl == nil {
		return r
	}
	r.endpoint = eurl
	r.endpoint.RawQuery = query
	return r
}

func (r *ClientExecutorRequest) SetQuery(params url.Values) *ClientExecutorRequest {
	r.endpoint.RawQuery = params.Encode()
	return r
}

func (r *ClientExecutorRequest) SetRequestHeader(key string, value string) *ClientExecutorRequest {
	r.headers[key] = value
	return r
}

func (r *ClientExecutorRequest) SetJSONBody(input interface{}) error {
	buf := bytes.NewBuffer(make([]byte, 0))
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(input); err != nil {
		return fmt.Errorf("error encoding json body: %s", err)
	}
	r.body = buf.Bytes()
	r.headers["Content-Type"] = "application/json"
	return nil
}

func (r *ClientExecutorRequest) Get() (*ClientExecutorResponse, error) {
	return r.do("GET")
}

func (r *ClientExecutorRequest) Post() (*ClientExecutorResponse, error) {
	return r.do("POST")
}

func (r *ClientExecutorRequest) Put() (*ClientExecutorResponse, error) {
	return r.do("PUT")
}

func (r *ClientExecutorRequest) Delete() (*ClientExecutorResponse, error) {
	return r.do("DELETE")
}

func (r *ClientExecutorRequest) do(method string) (*ClientExecutorResponse, error) {
	var body io.Reader
	if r.body != nil && len(r.body) > 0 {
		body = bytes.NewReader(r.body)
	}

	req, err := http.NewRequest(method, r.endpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	if r.debug {
		raw, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(raw))
	}

	res, err := r.client.Do(req)
	response := &ClientExecutorResponse{res: res}
	if r.debug && res != nil {
		raw, _ := httputil.DumpResponse(res, true)
		log.Println(string(raw))
	}
	if err != nil {
		return response, fmt.Errorf("error requesting %s %s: %s", req.Method, req.URL, err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return response, NewSumologicError(response)
	}

	return response, nil
}
