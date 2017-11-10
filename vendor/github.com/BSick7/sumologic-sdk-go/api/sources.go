package api

import (
	"fmt"
	"time"
)

// SumoLogic API Reference
// https://help.sumologic.com/APIs/01Collector-Management-API/Source-API
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
	ID                         int             `json:"id"`
	Name                       string          `json:"name"`
	SourceType                 string          `json:"sourceType,omitempty"`
	Description                string          `json:"description,omitempty"`
	Category                   string          `json:"category,omitempty"`
	HostName                   string          `json:"hostName,omitempty"`
	TimeZone                   string          `json:"timeZone,omitempty"`
	ForceTimeZone              bool            `json:"forceTimeZone,omitempty"`
	AutomaticDateParsing       bool            `json:"automaticDateParsing,omitempty"`
	MultilineProcessingEnabled bool            `json:"multilineProcessingEnabled,omitempty"`
	UseAutolineMatching        bool            `json:"useAutolineMatching,omitempty"`
	ManualPrefixRegexp         string          `json:"manualPrefixRegexp,omitempty"`
	MessagePerRequest          bool            `json:"messagePerRequest"`
	DefaultDateFormat          string          `json:"defaultDateFormat,omitempty"`
	DefaultDateFormats         []*DateFormat   `json:"defaultDateFormats,omitempty"`
	Filters                    []*SourceFilter `json:"filters,omitempty"`
	CutoffTimestamp            time.Time       `json:"-"`
	CutoffTimestampMs          int64           `json:"cutoffTimestamp,omitempty"`
	CutoffRelativeTime         string          `json:"cutoffRelativeTime,omitempty"`
	PathExpression             string          `json:"pathExpression,omitempty"`
	Blacklist                  []string        `json:"blacklist,omitempty"`
	Encoding                   string          `json:"encoding,omitempty"`
}

// This will coerce CutoffTimestampMs to CutoffTimestamp
func (s *Source) SyncTimestamp() {
	// Sumologic passes this around as number of milliseconds since epoch
	// time.Unix() returns number of seconds
	s.CutoffTimestampMs = s.CutoffTimestamp.Unix() * 1000
}

// This will coerce CutoffTimestamp to CutoffTimestampMs
func (s *Source) SyncTimestampMs() {
	s.CutoffTimestamp = time.Unix(s.CutoffTimestampMs*1000, 0)
}

type SourceCreate struct {
	SourceType        string  `json:"sourceType"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Category          string  `json:"category"`
	MessagePerRequest *bool   `json:"messagePerRequest,omitempty"`
	PathExpression    *string `json:"pathExpression,omitempty"`
}

type SourceFilter struct {
	FilterType string `json:"filterType"`
	Name       string `json:"name"`
	Regexp     string `json:"regexp"`
	Mask       string `json:"mask,omitempty"`
}

type DateFormat struct {
	Format  string `json:"format"`
	Locator string `json:"locator,omitempty"`
}

func (s *Sources) List() ([]*Source, error) {
	req, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources", s.collectorID))

	res, err := req.Get()
	if err != nil {
		return nil, err
	}

	type listResponse struct {
		Sources []*Source `json:"sources"`
	}
	list := &listResponse{}
	if err := res.BodyJSON(list); err != nil {
		return nil, err
	}
	if list.Sources != nil {
		for _, source := range list.Sources {
			source.SyncTimestamp()
		}
	}
	return list.Sources, nil
}

func (s *Sources) Get(id int) (*Source, error) {
	req, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, id))

	res, err := req.Get()
	if err != nil {
		return nil, err
	}

	type getResponse struct {
		Source *Source `json:"source"`
	}
	item := &getResponse{}
	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}
	if item.Source != nil {
		item.Source.SyncTimestamp()
	}
	return item.Source, nil
}

func (s *Sources) Exists(id int) (bool, error) {
	return IsObjectFound(s.Get(id))
}

func (s *Sources) Create(source *SourceCreate) (*Source, error) {
	req, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d/sources", s.collectorID))

	type postRequest struct {
		Source *SourceCreate `json:"source"`
	}
	req.SetJSONBody(&postRequest{Source: source})

	res, err := req.Post()
	if err != nil {
		return nil, err
	}

	type postResponse struct {
		Source *Source `json:"source"`
	}
	item := &postResponse{}
	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}
	return item.Source, nil
}

func (s *Sources) Update(source *Source) (*Source, error) {
	source.SyncTimestampMs()

	startreq, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	startreq.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, source.ID))

	startres, err := startreq.Get()
	if err != nil {
		return nil, err
	}

	finishreq, err := s.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	finishreq.SetEndpoint(fmt.Sprintf("/collectors/%d/sources/%d", s.collectorID, source.ID))
	finishreq.SetRequestHeader("If-Match", startres.Header("ETag"))

	type putRequest struct {
		Source *Source `json:"source"`
	}
	finishreq.SetJSONBody(&putRequest{Source: source})

	finishres, err := finishreq.Put()
	if err != nil {
		return nil, err
	}

	type putResponse struct {
		Source *Source `json:"source"`
	}
	item := &putResponse{}
	if err := finishres.BodyJSON(item); err != nil {
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

	if _, err := req.Delete(); err != nil {
		return err
	}

	return nil
}
