package labstack

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dghubble/sling"
	"github.com/labstack/gommon/log"
)

type (
	// Cube defines the LabStack cube service.
	Cube struct {
		sling          *sling.Sling
		requests       []*CubeRequest
		activeRequests int64
		started        int64
		mutex          *sync.RWMutex
		logger         *log.Logger

		// LabStack Account ID
		AccountID string

		// LabStack API key
		APIKey string

		// Number of requests in a batch
		BatchSize int

		// Interval in seconds to dispatch the batch
		DispatchInterval time.Duration

		// Additional tags
		Tags []string `json:"tags"`

		// TODO: To be implemented
		ClientLookup string
	}

	// CubeRequest defines a request payload to be recorded.
	CubeRequest struct {
		recovered bool
		ID        string    `json:"id"`
		Time      time.Time `json:"time"`
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
		Tags       []string `json:"tags"`
		Language   string   `json:"language"`
		Error      string   `json:"error"`
		StackTrace string   `json:"stack_trace"`
	}

	// CubeError defines the cube error.
	CubeError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func (c *Cube) resetRequests() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = make([]*CubeRequest, 0, c.BatchSize)
}

func (c *Cube) appendRequest(r *CubeRequest) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = append(c.requests, r)
}

func (c *Cube) listRequests() []*CubeRequest {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	requests := make([]*CubeRequest, len(c.requests))
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

// dispatch dispatches the requests batch.
func (c *Cube) dispatch() error {
	if len(c.requests) == 0 {
		return nil
	}

	ce := new(CubeError)
	_, err := c.sling.Post("").BodyJSON(c.listRequests()).Receive(nil, ce)
	if err != nil {
		return err
	}
	if ce.Code != 0 {
		return ce
	}

	return nil
}

// Start starts recording an HTTP request.
func (c *Cube) Start(r *http.Request, w http.ResponseWriter) (request *CubeRequest) {
	if c.started == 0 {
		go func() {
			d := time.Duration(c.DispatchInterval) * time.Second
			for range time.Tick(d) {
				c.dispatch()
			}
		}()
		atomic.AddInt64(&c.started, 1)
	}

	request = &CubeRequest{
		ID:        RequestID(r, w),
		Time:      time.Now(),
		Host:      r.Host,
		Path:      r.URL.Path,
		Method:    r.Method,
		UserAgent: r.UserAgent(),
		RemoteIP:  RealIP(r),
		Language:  "Go",
		Tags:      c.Tags,
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

// Recover handles a panic
func (c *Cube) Recover(r interface{}, cr *CubeRequest) {
	if r == nil {
		return
	}
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}
	stack := make([]byte, 4<<10) // 4 KB
	length := runtime.Stack(stack, false)
	cr.Error = err.Error()
	cr.StackTrace = string(stack[:length])
	cr.recovered = true
}

// Stop stops recording an HTTP request.
func (c *Cube) Stop(r *CubeRequest, status int, size int64) {
	if r.recovered {
		status = http.StatusInternalServerError
	}
	atomic.AddInt64(&c.activeRequests, -1)
	r.Status = status
	r.BytesOut = size
	r.Latency = int64(time.Now().Sub(r.Time))

	// Dispatch batch
	if r.Status >= 500 && r.Status < 600 || c.requestsLength() >= c.BatchSize {
		go func() {
			if err := c.dispatch(); err != nil {
				c.logger.Error(err)
			}
			c.resetRequests()
		}()
	}
}

func (e *CubeError) Error() string {
	return fmt.Sprintf("cube error, code=%d, message=%s", e.Code, e.Message)
}
