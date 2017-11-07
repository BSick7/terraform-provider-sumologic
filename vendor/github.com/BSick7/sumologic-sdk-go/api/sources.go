package api

import (
	"fmt"
)

type Sources struct {
	executor    *ClientExecutor
	collectorID int
}

func NewSources(executor *ClientExecutor, collectorID int) *Sources {
	return &Sources{
		executor:    executor,
		collectorID: collectorID,
	}
}

type Source struct {
	ID                         int               `json:"id"`
	Name                       string            `json:"name"`
	Description                string            `json:"description,omitempty"`
	Category                   string            `json:"category"`
	HostName                   string            `json:"hostName,omitempty"`
	AutomaticDateParsing       bool              `json:"automaticDateParsing"`
	MultilineProcessingEnabled bool              `json:"multilineProcessingEnabled"`
	UseAutolineMatching        bool              `json:"useAutolineMatching"`
	ForceTimeZone              bool              `json:"forceTimeZone"`
	TimeZone                   string            `json:"timeZone"`
	Filters                    []SourceFilter    `json:"filters"`
	CutoffTimestamp            int               `json:"cutoffTimestamp"`
	Encoding                   string            `json:"encoding"`
	MessagePerRequest          bool              `json:"messagePerRequest"`
	PathExpression             string            `json:"pathExpression,omitempty"`
	Blacklist                  []SourceBlacklist `json:"blacklist,omitempty"`
	SourceType                 string            `json:"sourceType"`
	Metrics                    []string          `json:"metrics,omitempty"`
	Interval                   int               `json:"interval,omitempty"`
	Alive                      bool              `json:"alive,omitempty"`
}

type SourceFilter struct {
	FilterType string `json:"filterType"`
	Name       string `json:"name"`
	Regexp     string `json:"regexp"`
}

type SourceBlacklist struct {
}

func (s *Sources) List() ([]*Source, error) {
	req, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources", s.collectorID))

	if err := req.Get(); err != nil {
		return nil, err
	}

	type listResponse struct {
		Sources []*Source `json:"sources"`
	}
	list := &listResponse{}
	if err := req.GetJSONBody(list); err != nil {
		return nil, err
	}
	return list.Sources, nil
}

func (s *Sources) Get(id int) (*Source, error) {
	req, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, id))

	if err := req.Get(); err != nil {
		return nil, err
	}

	type getResponse struct {
		Source *Source `json:"source"`
	}
	item := &getResponse{}
	if err := req.GetJSONBody(item); err != nil {
		return nil, err
	}
	return item.Source, nil
}

func (s *Sources) Create(source *Source) (*Source, error) {
	req, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources", s.collectorID))

	type postRequest struct {
		Source *Source `json:"source"`
	}
	req.SetJSONBody(&postRequest{Source: source})

	if err := req.Post(); err != nil {
		return nil, err
	}

	type postResponse struct {
		Source *Source `json:"source"`
	}
	item := &postResponse{}
	if err := req.GetJSONBody(item); err != nil {
		return nil, err
	}
	return item.Source, nil
}

func (s *Sources) Update(source *Source) (*Source, error) {
	startreq, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	startreq.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, source.ID))

	if err := startreq.Put(); err != nil {
		return nil, err
	}

	finishreq, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	finishreq.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, source.ID))
	finishreq.SetRequestHeader("If-Match", startreq.GetResponseHeader("ETag"))

	type putRequest struct {
		Source *Source `json:"source"`
	}
	finishreq.SetJSONBody(&putRequest{Source: source})

	if err := finishreq.Put(); err != nil {
		return nil, err
	}

	type putResponse struct {
		Source *Source `json:"source"`
	}
	item := &putResponse{}
	if err := finishreq.GetJSONBody(item); err != nil {
		return nil, err
	}
	return item.Source, nil
}

func (s *Sources) Delete(source *Source) error {
	req, err := s.executor.NewRequest()
	if err != nil {
		return err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, source.ID))

	if err := req.Delete(); err != nil {
		return err
	}

	return nil
}
