package cube

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"sync/atomic"

	"encoding/json"

	"fmt"

	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type (
	// Config defines the config for Cube middleware.
	Config struct {
		Skipper       Skipper
		Node          string        `json:"node"`
		Group         string        `json:"group"`
		APIKey        string        `json:"api_key"`
		BatchSize     int           `json:"batch_size"`
		BatchInterval time.Duration `json:"batch_interval"`
		ClientLookup  string        `json:"client_lookup"`
	}

	cube struct {
		client         *http.Client
		requests       []*request
		activeRequests int64
		mutex          sync.RWMutex
		logger         *log.Logger
		Config
	}

	request struct {
		ID             string    `json:"id"`
		Time           time.Time `json:"time"`
		Node           string    `json:"node"`
		Group          string    `json:"group"`
		Host           string    `json:"host"`
		Path           string    `json:"path"`
		Method         string    `json:"method"`
		Status         int       `json:"status"`
		BytesIn        int64     `json:"bytes_in"`
		BytesOut       int64     `json:"bytes_out"`
		Latency        int64     `json:"latency"`
		ClientID       string    `json:"client_id"`
		RemoteIP       string    `json:"remote_ip"`
		UserAgent      string    `json:"user_agent"`
		ActiveRequests int64     `json:"active_requests"`
		// TODO: CPU, Uptime, Memory
	}

	// Skipper defines a function to conditionally skip the middleware
	Skipper func(r *http.Request) bool
)

const (
	apiURL = "https://api.labstack.com/cube"
)

func (c *cube) resetRequests() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = make([]*request, 0, c.BatchSize)
}

func (c *cube) appendRequest(r *request) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requests = append(c.requests, r)
}

func (c *cube) send() (err error) {
	if len(c.requests) == 0 {
		return
	}

	c.mutex.RLock()
	buf := new(bytes.Buffer)
	if err = json.NewEncoder(buf).Encode(c.requests); err != nil {
		return
	}
	c.mutex.RUnlock()
	req, err := http.NewRequest(echo.POST, apiURL, buf)
	if err != nil {
		return
	}
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+c.APIKey)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	res, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			body = []byte(err.Error())
		}
		return fmt.Errorf("cube: sending requests error=%s", body)
	}

	c.resetRequests()
	return
}

// MiddlewareEcho implements Cube middleware for Echo.
func MiddlewareEcho(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.BatchSize == 0 {
		config.BatchSize = 60
	}
	if config.BatchInterval == 0 {
		config.BatchInterval = 60
	}
	if config.Skipper == nil {
		config.Skipper = func(*http.Request) bool {
			return false
		}
	}

	// Initialize
	c := &cube{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		logger: log.New("cube"),
		Config: config,
	}
	c.resetRequests()
	go func() {
		for range time.Tick(config.BatchInterval * time.Second) {
			c.send()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if config.Skipper(ctx.Request()) {
				return next(ctx)
			}

			req := ctx.Request()
			res := ctx.Response()
			start := time.Now()
			r := &request{
				Time:      time.Now(),
				Node:      config.Node,
				Group:     config.Group,
				Host:      req.Host,
				Path:      req.URL.Path,
				Method:    req.Method,
				UserAgent: req.UserAgent(),
				RemoteIP:  ctx.RealIP(),
				ClientID:  ctx.RealIP(),
			}
			atomic.AddInt64(&c.activeRequests, 1)
			r.ActiveRequests = c.activeRequests
			c.appendRequest(r)

			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			atomic.AddInt64(&c.activeRequests, -1)
			stop := time.Now()
			r.Status = res.Status
			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}
			l, err := strconv.ParseInt(cl, 10, 64)
			if err != nil {
				c.logger.Error(err)
			}
			r.BytesIn = l
			r.BytesOut = res.Size
			l = int64(stop.Sub(start))
			r.Latency = int64(stop.Sub(start))

			// Send batch
			if len(c.requests) == config.BatchSize {
				go func() {
					if err := c.send(); err != nil {
						c.logger.Error(err)
					}
				}()
			}

			return
		}
	}
}
