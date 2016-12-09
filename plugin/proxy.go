package plugin

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"

	"github.com/labstack/echo"
)

type (
	Proxy struct {
		Base     `json:",squash"`
		Balance  string    `json:"balance"`
		Targets  []*Target `json:"targets"`
		balancer Balancer
	}

	Target struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url"`
		url  *url.URL
	}

	Random struct {
		targets []*Target
		random  *rand.Rand
	}

	RoundRobin struct {
		targets []*Target
		i       uint32
	}

	Balancer interface {
		Next(int) *Target
		Length() int
	}
)

func (p *Proxy) Init() (err error) {
	for _, t := range p.Targets {
		t.url, err = url.Parse(t.URL)
		if err != nil {
			return
		}
	}

	// Balancer
	switch p.Balance {
	case "round-robin":
		p.balancer = &RoundRobin{targets: p.Targets}
	default: // Random
		p.balancer = &Random{
			targets: p.Targets,
			random:  rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		}
	}

	return
}

func (*Proxy) Priority() int {
	return 1
}

func (p *Proxy) Process(echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		i := -1
		t := p.balancer.Next(i).url

		// Next:
		// 	*outReq.URL = *p.balancer.Next(i).url // Shallow copy
		// 	path := c.P(0)
		// 	if path != "" {
		// 		if path[0] != '/' { // Prepend '/' if necessary
		// 			path = "/" + path
		// 		}
		// 		outReq.URL.Path = path
		// 		outReq.URL.RawQuery = c.QueryString()
		// 	}
		// 	outReq.Header = req.Header
		// 	p.Logger.Infof("proxy: out request, url=%v", outReq.URL)

		if req.Header.Get(echo.HeaderUpgrade) == "websocket" {
			p.wsProxy(t).ServeHTTP(res, req)
		} else {
			p.httpProxy(t).ServeHTTP(res, req)
		}

		// 	// Out request
		// 	outReq := &http.Request{
		// 		Method:        req.Method,
		// 		Close:         false, // Persistent connection
		// 		URL:           new(url.URL),
		// 		Body:          req.Body,
		// 		ContentLength: int64(req.ContentLength),
		// 	}
		// Next:
		// 	*outReq.URL = *p.balancer.Next(i).url // Shallow copy
		// 	path := c.P(0)
		// 	if path != "" {
		// 		if path[0] != '/' { // Prepend '/' if necessary
		// 			path = "/" + path
		// 		}
		// 		outReq.URL.Path = path
		// 		outReq.URL.RawQuery = c.QueryString()
		// 	}
		// 	outReq.Header = req.Header
		// 	p.Logger.Infof("proxy: out request, url=%v", outReq.URL)
		//
		// 	// Out response
		// 	outRes, err := p.transport.RoundTrip(outReq)
		// 	if err != nil {
		// 		p.Logger.Errorf("proxy: error=%s", err)
		// 		i++
		// 		if i == p.balancer.Length() {
		// 			return echo.NewHTTPError(http.StatusBadGateway)
		// 		}
		// 		p.Logger.Warnf("proxy: trying target=%s", outReq.URL)
		// 		goto Next
		// 	}
		// 	p.Logger.Infof("proxy: out response, status=%d", outRes.StatusCode)
		// 	defer outRes.Body.Close()
		// 	for k := range outRes.Header { // Copy headers
		// 		res.Header().Add(k, outRes.Header.Get(k))
		// 	}
		// 	res.WriteHeader(outRes.StatusCode)
		// 	_, err = io.Copy(res, outRes.Body)
		return
	}
}

func (r *Random) Next(i int) *Target {
	if i == -1 {
		return r.targets[r.random.Intn(len(r.targets))]
	}
	return r.targets[i]
}

func (r *Random) Length() int {
	return len(r.targets)
}

func (r *RoundRobin) Next(i int) *Target {
	if i == -1 {
		if r.i == uint32(len(r.targets)) {
			r.i = 0
		}
		atomic.AddUint32(&r.i, 1)
		return r.targets[r.i]
	}
	return r.targets[i]
}

func (r *RoundRobin) Length() int {
	return len(r.targets)
}

func (p *Proxy) httpProxy(u *url.URL) http.Handler {
	return httputil.NewSingleHostReverseProxy(u)
}

func (p *Proxy) wsProxy(u *url.URL) http.Handler {
	return websocket.Handler(func(in *websocket.Conn) {
		defer in.Close()

		r := in.Request()
		t := "ws://" + u.Host + r.RequestURI
		out, err := websocket.Dial(t, "", r.Header.Get("Origin"))
		if err != nil {
			p.Logger.Errorf("ws proxy error, target=%s, err=%v", t, err)
			return
		}
		defer out.Close()

		errc := make(chan error, 2)
		cp := func(w io.Writer, r io.Reader) {
			_, err := io.Copy(w, r)
			errc <- err
		}

		go cp(in, out)
		go cp(out, in)
		err = <-errc
		if err != nil && err != io.EOF {
			p.Logger.Errorf("ws proxy error, url=%s, err=%v", r.URL, err)
		}
	})
}
