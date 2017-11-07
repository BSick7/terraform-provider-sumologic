package api

import (
	"fmt"
	"net/url"
	"strconv"
)

type Collectors struct {
	executor *ClientExecutor
}

func NewCollectors(executor *ClientExecutor) *Collectors {
	return &Collectors{
		executor: executor,
	}
}

type Collector struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	TimeZone           string          `json:"timeZone,omitempty"`
	Links              []CollectorLink `json:"links,omitempty"`
	Ephemeral          bool            `json:"ephemeral,omitempty"`
	SourceSyncMode     string          `json:"sourceSyncMode,omitempty"`
	CollectorType      string          `json:"collectorType"`
	CollectorVersion   string          `json:"collectorVersion,omitempty"`
	Description        string          `json:"description,omitempty"`
	OsArch             string          `json:"osArch,omitempty"`
	OsVersion          string          `json:"osVersion,omitempty"`
	OsName             string          `json:"osName,omitempty"`
	OsTime             int64           `json:"osTime,omitempty"`
	Category           string          `json:"category"`
	LastSeenAlive      int64           `json:"lastSeenAlive,omitempty"`
	Alive              bool            `json:"alive,omitempty"`
	CutoffTimestamp    int64           `json:"cutoffTimestamp,omitempty"`
	CutoffRelativeTime string          `json:"cutoffRelativeTime,omitempty"`
	TargetCPU          int64           `json:"targetCPU,omitempty"`
	HostName           string          `json:"hostName,omitempty"`
}

type CollectorCreate struct {
	CollectorType string `json:"collectorType"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Category      string `json:"category"`
}

type CollectorLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func (c *Collectors) Sources(collectorID int) *Sources {
	return NewSources(c.executor, collectorID)
}

func (c *Collectors) List(offset int, limit int) ([]*Collector, error) {
	req, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint("/collectors")
	req.SetQuery(url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	})

	if err := req.Get(); err != nil {
		return nil, err
	}

	type listResponse struct {
		Collectors []*Collector `json:"collectors"`
	}
	list := &listResponse{}
	if err := req.GetJSONBody(list); err != nil {
		return nil, err
	}
	return list.Collectors, nil
}

func (c *Collectors) Get(id int) (*Collector, error) {
	req, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d", id))

	if err := req.Get(); err != nil {
		return nil, err
	}

	type getResponse struct {
		Collector *Collector `json:"collector"`
	}
	item := &getResponse{}
	if err := req.GetJSONBody(item); err != nil {
		return nil, err
	}
	return item.Collector, nil
}

func (c *Collectors) Create(collector *CollectorCreate) (*Collector, error) {
	req, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint("/collectors")

	type postRequest struct {
		Collector *CollectorCreate `json:"collector"`
	}
	req.SetJSONBody(&postRequest{Collector: collector})

	if err := req.Post(); err != nil {
		return nil, err
	}

	type postResponse struct {
		Collector *Collector `json:"collector"`
	}
	item := &postResponse{}
	if err := req.GetJSONBody(item); err != nil {
		return nil, err
	}
	return item.Collector, nil
}

func (c *Collectors) Update(collector *Collector) (*Collector, error) {
	startreq, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	startreq.SetEndpoint(fmt.Sprintf("/collectors/%d", collector.ID))

	if err := startreq.Put(); err != nil {
		return nil, err
	}

	finishreq, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	finishreq.SetEndpoint(fmt.Sprintf("/collectors/%d", collector.ID))
	finishreq.SetRequestHeader("If-Match", startreq.GetResponseHeader("ETag"))

	type putRequest struct {
		Collector *Collector `json:"collector"`
	}
	finishreq.SetJSONBody(&putRequest{Collector: collector})

	if err := finishreq.Put(); err != nil {
		return nil, err
	}

	type putResponse struct {
		Collector *Collector `json:"collector"`
	}
	item := &putResponse{}
	if err := finishreq.GetJSONBody(item); err != nil {
		return nil, err
	}
	return item.Collector, nil
}

func (c *Collectors) Delete(collector *Collector) error {
	req, err := c.executor.NewRequest()
	if err != nil {
		return err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d", collector.ID))

	if err := req.Delete(); err != nil {
		return err
	}

	return nil
}
