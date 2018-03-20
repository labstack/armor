package cube

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty"
	"github.com/labstack/gommon/log"
	"github.com/shirou/gopsutil/process"
)

type (
	Cube struct {
		Options
		uptime         uint64
		cpu            float32
		memory         uint64
		requests       []*Request
		activeRequests int64
		mutex          sync.RWMutex
		client         *resty.Client
		logger         *log.Logger
	}

	Options struct {
		// Node id
		Node string

		// Tags
		Tags []string

		// Number of requests in a batch
		BatchSize int

		// Interval in seconds to dispatch the batch
		DispatchInterval time.Duration
	}

	Request struct {
		ID         string   `json:"id,omitempty"`
		Time       int64    `json:"time,omitempty"`
		Tags       []string `json:"tags,omitempty"`
		Host       string   `json:"host,omitempty"`
		Path       string   `json:"path,omitempty"`
		Method     string   `json:"method,omitempty"`
		Status     int      `json:"status,omitempty"`
		BytesIn    int64    `json:"bytes_in,omitempty"`
		BytesOut   int64    `json:"bytes_out,omitempty"`
		Latency    int64    `json:"latency,omitempty"`
		ClientID   string   `json:"client_id,omitempty"`
		RemoteIP   string   `json:"remote_ip,omitempty"`
		UserAgent  string   `json:"user_agent,omitempty"`
		Active     int64    `json:"active,omitempty"`
		Error      string   `json:"error,omitempty"`
		StackTrace string   `json:"stack_trace,omitempty"`
		Node       string   `json:"node,omitempty"`
		Uptime     uint64   `json:"uptime,omitempty"`
		CPU        float32  `json:"cpu,omitempty"`
		Memory     uint64   `json:"memory,omitempty"`
	}

	// APIError struct {
	// 	Code    int    `json:"code"`
	// 	Message string `json:"message"`
	// }
)

func New(apiKey string, options Options) *Cube {
	c := new(Cube)
	c.Options = options
	c.client = resty.New().
		SetHostURL("https://api.labstack.com").
		SetAuthToken(apiKey).
		SetHeader("User-Agent", "labstack/cube")

	// Defaults
	if c.Node == "" {
		c.Node, _ = os.Hostname()
	}
	if c.BatchSize == 0 {
		c.BatchSize = 60
	}
	if c.DispatchInterval == 0 {
		c.DispatchInterval = 60
	}

	// Dispatch daemon
	go func() {
		d := time.Duration(c.DispatchInterval) * time.Second
		for range time.Tick(d) {
			c.Dispatch()
		}
	}()

	// System daemon
	go func() {
		p, _ := process.NewProcess(int32(os.Getpid()))
		t, _ := p.CreateTime()

		for range time.Tick(time.Second) {
			c.uptime = uint64(time.Now().Unix() - t/1000)
			cpu, _ := p.CPUPercent()
			c.cpu = float32(cpu)
			mem, _ := p.MemoryInfo()
			c.memory = mem.RSS
		}
	}()

	return c
}

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

// Dispatch dispatches the requests batch.
func (c *Cube) Dispatch() {
	if len(c.requests) == 0 {
		return
	}

	// err := new(APIError)
	res, err := c.client.R().
		SetBody(c.listRequests()).
		// SetError(err).
		Post("/cube")
	if err != nil {
		c.logger.Error(err)
		return
	}
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		c.logger.Error(res.Body())
	}

	c.resetRequests()
}

// Start starts cording an HTTP request.
func (c *Cube) Start(r *Request) {
	atomic.AddInt64(&c.activeRequests, 1)

	r.Time = time.Now().UnixNano() / 1000
	r.Active = c.activeRequests
	r.Node = c.Node
	r.Uptime = c.uptime
	r.CPU = c.cpu
	r.Memory = c.memory
	r.Tags = c.Tags

	c.appendRequest(r)
}

// Stop stops recording an HTTP request.
func (c *Cube) Stop(r *Request) {
	atomic.AddInt64(&c.activeRequests, -1)

	r.Latency = time.Now().UnixNano()/1000 - r.Time

	// Dispatch batch
	if c.requestsLength() >= c.BatchSize {
		go func() {
			c.Dispatch()
		}()
	}
}

// RequestID returns the request ID from the request or response.
// func RequestID(r *http.Request, w http.ResponseWriter) string {
// 	id := r.Header.Get("X-Request-ID")
// 	if id == "" {
// 		id = w.Header().Get("X-Request-ID")
// 	}
// 	return id
// }

// RealIP returns the real IP from the request.
// func RealIP(r *http.Request) string {
// 	ra := r.RemoteAddr
// 	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
// 		ra = strings.Split(ip, ", ")[0]
// 	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
// 		ra = ip
// 	} else {
// 		ra, _, _ = net.SplitHostPort(ra)
// 	}
// 	return ra
// }

// func (e *APIError) Error() string {
// 	return e.Message
// }
