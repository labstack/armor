package cube

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"sync/atomic"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type (
	Cube struct {
		logger         *log.Logger
		mutex          sync.RWMutex
		activeRequests int64
		requests       []*Request
		Skipper        Skipper
		CacheLimit     int
	}

	Request struct {
		Time      int64  `json:"time"`
		Path      string `json:"path"`
		Method    string `json:"method"`
		Active    int64  `json:"active"`
		Status    int    `json:"status"`
		BytesIn   int64  `json:"bytes_in"`
		BytesOut  int64  `json:"bytes_out"`
		Latency   int64  `json:"latency"`
		ClientID  string `json:"client_id"`
		RemoteIP  string `json:"remote_ip"`
		UserAgent string `json:"user_agent"`
	}

	Skipper func(r *http.Request) bool
)

func New() *Cube {
	return &Cube{
		Skipper: func(*http.Request) bool {
			return false
		},
	}
}

func (c *Cube) reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = nil
}

func (c *Cube) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		if c.Skipper(ctx.Request()) {
			return next(ctx)
		}

		req := ctx.Request()
		res := ctx.Response()
		start := time.Now()
		r := &Request{
			Time: time.Now().UnixNano(),
		}
		atomic.AddInt64(&c.activeRequests, 1)
		if err = next(ctx); err != nil {
			ctx.Error(err)
		}
		atomic.AddInt64(&c.activeRequests, -1)
		stop := time.Now()
		r.Path = req.URL.Path
		r.Method = req.Method
		r.Status = res.Status
		cl := req.Header.Get(echo.HeaderContentLength)
		if cl == "" {
			cl = "0"
		}
		l, err := strconv.ParseInt(cl, 10, 64)
		if err != nil {
			ctx.Error(err)
		}
		r.BytesIn = l
		r.BytesOut = res.Size
		l = int64(stop.Sub(start))
		r.Latency = int64(stop.Sub(start))
		r.UserAgent = req.UserAgent()
		r.RemoteIP = ctx.RealIP()
		r.ClientID = ctx.RealIP()

		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.requests = append(c.requests, r)

		return
	}
}

func (c *Cube) Requests() []*Request {
	c.mutex.RLock()
	requests := make([]*Request, len(c.requests))
	for i, r := range c.requests {
		requests[i] = r
	}
	c.mutex.RUnlock()
	c.reset()
	return requests
}
