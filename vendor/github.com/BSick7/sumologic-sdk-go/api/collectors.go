package api

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// SumoLogic API Reference
// https://help.sumologic.com/APIs/01Collector-Management-API/Collector-API-Methods-and-Examples
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
	CutoffTimestamp    time.Time       `json:"-"`
	CutoffTimestampMs  int64           `json:"cutoffTimestamp,omitempty"`
	CutoffRelativeTime string          `json:"cutoffRelativeTime,omitempty"`
	TargetCPU          int64           `json:"targetCPU,omitempty"`
	HostName           string          `json:"hostName,omitempty"`
}

// This will coerce CutoffTimestampMs to CutoffTimestamp
func (c *Collector) SyncTimestamp() {
	// Sumologic passes this around as number of milliseconds since epoch
	// time.Unix() returns number of seconds
	c.CutoffTimestampMs = c.CutoffTimestamp.Unix() * 1000
}

// This will coerce CutoffTimestamp to CutoffTimestampMs
func (c *Collector) SyncTimestampMs() {
	c.CutoffTimestamp = time.Unix(c.CutoffTimestampMs*1000, 0)
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

	res, err := req.Get()
	if err != nil {
		return nil, err
	}

	type listResponse struct {
		Collectors []*Collector `json:"collectors"`
	}
	list := &listResponse{}
	if err := res.BodyJSON(list); err != nil {
		return nil, err
	}
	if list.Collectors != nil {
		for _, collector := range list.Collectors {
			collector.SyncTimestamp()
		}
	}
	return list.Collectors, nil
}

func (c *Collectors) Get(id int) (*Collector, error) {
	req, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/collectors/%d", id))

	res, err := req.Get()
	if err != nil {
		return nil, err
	}

	type getResponse struct {
		Collector *Collector `json:"collector"`
	}
	item := &getResponse{}
	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}
	if item.Collector != nil {
		item.Collector.SyncTimestamp()
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

	res, err := req.Post()
	if err != nil {
		return nil, err
	}

	type postResponse struct {
		Collector *Collector `json:"collector"`
	}
	item := &postResponse{}
	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}
	return item.Collector, nil
}

func (c *Collectors) Update(collector *Collector) (*Collector, error) {
	collector.SyncTimestampMs()

	startreq, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	startreq.SetEndpoint(fmt.Sprintf("/collectors/%d", collector.ID))

	startres, err := startreq.Get()
	if err != nil {
		return nil, err
	}

	finishreq, err := c.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	finishreq.SetEndpoint(fmt.Sprintf("/collectors/%d", collector.ID))
	finishreq.SetRequestHeader("If-Match", startres.Header("ETag"))

	type putRequest struct {
		Collector *Collector `json:"collector"`
	}
	finishreq.SetJSONBody(&putRequest{Collector: collector})

	finishres, err := finishreq.Put()
	if err != nil {
		return nil, err
	}

	type putResponse struct {
		Collector *Collector `json:"collector"`
	}
	item := &putResponse{}
	if err := finishres.BodyJSON(item); err != nil {
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

	if _, err := req.Delete(); err != nil {
		return err
	}

	return nil
}
