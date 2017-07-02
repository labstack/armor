package labstack

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dghubble/sling"
	"github.com/labstack/gommon/log"
)

type (
	// Cube defines the Cube service.
	Cube struct {
		sling          *sling.Sling
		requests       []*Request
		activeRequests int64
		mutex          sync.RWMutex
		logger         *log.Logger

		// Node name
		Node string

		// Node group
		Group string

		// LabStack API key
		APIKey string

		// Number of requests in a batch
		BatchSize int

		// Interval in seconds to dispatch the batch
		DispatchInterval time.Duration

		// TODO: To be implemented
		ClientLookup string
	}

	// Request defines a request payload to be recorded.
	Request struct {
		ID        string    `json:"id"`
		Time      time.Time `json:"time"`
		Node      string    `json:"node"`
		Group     string    `json:"group"`
		Host      string    `json:"host"`
		Path      string    `json:"path"`
		Method    string    `json:"method"`
		Status    int       `json:"status"`
		BytesIn   int64     `json:"bytes_in"`
		BytesOut  int64     `json:"bytes_out"`
		Latency   int64     `json:"latency"`
		ClientID  string    `json:"client_id"`
		RemoteIP  string    `json:"remote_ip"`
		UserAgent string    `json:"user_agent"`
		Active    int64     `json:"active"`
		// TODO: CPU, Uptime, Memory
	}
)

func (c *Cube) resetRequests() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = make([]*Request, 0, c.BatchSize)
}

func (c *Cube) appendRequest(r *Request) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = append(c.requests, r)
}

func (c *Cube) listRequests() []*Request {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	requests := make([]*Request, len(c.requests))
	for i, r := range c.requests {
		requests[i] = r
	}
	return requests
}

func (c *Cube) requestsLength() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.requests)
}

func (c *Cube) dispatch() (err error) {
	if len(c.requests) == 0 {
		return
	}
	res, err := c.sling.Post("/cube").BodyJSON(c.listRequests()).Receive(nil, nil)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("cube: requests dispatching error=%s", err)
	}
	return
}

// Cube returns the Cube service.
func (c *Client) Cube() (cube *Cube) {
	cube = &Cube{
		sling:            c.sling.Path("/cube"),
		logger:           log.New("cube"),
		BatchSize:        60,
		DispatchInterval: 60,
	}
	cube.resetRequests()
	go func() {
		d := time.Duration(cube.DispatchInterval) * time.Second
		for range time.Tick(d) {
			cube.dispatch()
		}
	}()
	return
}

// Start starts recording an HTTP request.
func (c *Cube) Start(r *http.Request, w http.ResponseWriter) (request *Request) {
	request = &Request{
		Time:      time.Now(),
		Node:      c.Node,
		Group:     c.Group,
		Host:      r.Host,
		Path:      r.URL.Path,
		Method:    r.Method,
		UserAgent: r.UserAgent(),
		RemoteIP:  realIP(r),
	}
	request.ClientID = request.RemoteIP
	atomic.AddInt64(&c.activeRequests, 1)
	request.Active = c.activeRequests
	cl := r.Header.Get("Content-Length")
	if cl == "" {
		cl = "0"
	}
	l, err := strconv.ParseInt(cl, 10, 64)
	if err != nil {
		c.logger.Error(err)
	}
	request.BytesIn = l
	c.appendRequest(request)
	return
}

// Stop stops recording an HTTP request.
func (c *Cube) Stop(request *Request, status int, size int64) {
	atomic.AddInt64(&c.activeRequests, -1)
	request.Status = status
	request.BytesOut = size
	request.Latency = int64(time.Now().Sub(request.Time))

	// Dispatch batch
	if c.requestsLength() >= c.BatchSize {
		go func() {
			if err := c.dispatch(); err != nil {
				c.logger.Error(err)
			}
			c.resetRequests()
		}()
	}
}
