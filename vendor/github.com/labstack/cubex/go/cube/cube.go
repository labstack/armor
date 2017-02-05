package cube

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"sync/atomic"

	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
)

type (
	Cube struct {
		uptime        time.Time
		request       int64
		activeRequest int64
		bytesIn       int64
		bytesOut      int64
		latency       *Window
		endpoint      map[string]int64
		userAgent     map[string]int64
		remoteIP      map[string]int64
		status        map[int32]int64
		windowSize    int
		mutex         sync.Mutex
		Skipper       Skipper
		Tags          map[string]string `json:"tags"`
	}

	Skipper func(r *http.Request) bool
)

func New() *Cube {
	return &Cube{
		uptime:    time.Now(),
		endpoint:  map[string]int64{},
		userAgent: map[string]int64{},
		remoteIP:  map[string]int64{},
		status:    map[int32]int64{},
		latency:   NewWindow(50),
		Skipper: func(*http.Request) bool {
			return false
		},
		windowSize: 50,
		Tags:       map[string]string{},
	}
}

func (c *Cube) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		if c.Skipper(ctx.Request()) {
			return next(ctx)
		}

		req := ctx.Request()
		res := ctx.Response()
		start := time.Now()
		atomic.AddInt64(&c.activeRequest, 1)
		if err = next(ctx); err != nil {
			ctx.Error(err)
		}
		atomic.AddInt64(&c.activeRequest, -1)
		stop := time.Now()

		// Update (Acquire lock post request to prevent a deadlock)
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.request++
		cl := req.Header.Get(echo.HeaderContentLength)
		if cl == "" {
			cl = "0"
		}
		l, err := strconv.ParseInt(cl, 10, 64)
		if err != nil {
			ctx.Error(err)
		}
		c.bytesIn += l
		c.bytesOut += res.Size
		l = int64(stop.Sub(start))
		c.latency.Push(l)
		c.endpoint[req.URL.Path]++
		ua, _ := user_agent.New(req.UserAgent()).Browser()
		c.userAgent[ua]++
		c.remoteIP[ctx.RealIP()]++
		c.status[int32(res.Status)]++

		return
	}
}

func (c *Cube) Data() (d *Data) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	d = &Data{
		Uptime:         int64(time.Now().Sub(c.uptime).Seconds()),
		Request:        c.request,
		ActiveRequest:  c.activeRequest,
		BytesIn:        c.bytesIn,
		BytesOut:       c.bytesOut,
		AverageLatency: c.latency.Mean(),
		Endpoint:       c.endpoint,
		UserAgent:      c.userAgent,
		RemoteIP:       c.remoteIP,
		Status:         c.status,
		Tags:           c.Tags,
	}

	// Reset data
	c.reset()

	return
}

func (c *Cube) reset() {
	c.request = 0
	c.bytesIn = 0
	c.bytesOut = 0
	c.latency = NewWindow(c.windowSize)
	c.endpoint = map[string]int64{}
	c.userAgent = map[string]int64{}
	c.remoteIP = map[string]int64{}
	c.status = map[int32]int64{}
}
