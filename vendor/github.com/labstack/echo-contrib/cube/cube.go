package cube

import (
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/labstack-go"
)

type (
	// Config defines the config for Cube middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Node name
		Node string `json:"node"`

		// Node group
		Group string `json:"group"`

		// LabStack API key
		APIKey string `json:"api_key"`

		// Number of requests in a batch
		BatchSize int `json:"batch_size"`

		// Interval in seconds to dispatch the batch
		DispatchInterval time.Duration `json:"dispatch_interval"`

		// TODO: To be implemented
		ClientLookup string `json:"client_lookup"`
	}
)

var (
	// DefaultConfig is the default Cube middleware config.
	DefaultConfig = Config{
		Skipper:          middleware.DefaultSkipper,
		BatchSize:        60,
		DispatchInterval: 60,
	}
)

// Middleware implements Cube middleware.
func Middleware(apiKey string) echo.MiddlewareFunc {
	c := DefaultConfig
	c.APIKey = apiKey
	return MiddlewareWithConfig(c)
}

// MiddlewareWithConfig returns a Cube middleware with config.
// See: `Middleware()`.
func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.APIKey == "" {
		panic("echo: cube middleware requires an api key")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	if config.BatchSize == 0 {
		config.BatchSize = DefaultConfig.BatchSize
	}
	if config.DispatchInterval == 0 {
		config.DispatchInterval = DefaultConfig.DispatchInterval
	}

	// Initialize
	cube := labstack.NewClient(config.APIKey).Cube()
	cube.Node = config.Node
	cube.Group = config.Group
	cube.APIKey = config.APIKey
	cube.BatchSize = config.BatchSize
	cube.DispatchInterval = config.DispatchInterval
	cube.ClientLookup = config.ClientLookup

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}
			request := cube.Start(c.Request(), c.Response())
			if err = next(c); err != nil {
				c.Error(err)
			}
			cube.Stop(request, c.Response().Status, c.Response().Size)
			return
		}
	}
}
