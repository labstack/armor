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
		labstack.Cube

		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Number of requests in a batch
		BatchSize int

		// Interval in seconds to dispatch the batch
		DispatchInterval time.Duration
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
func Middleware(accountID, apiKey string) echo.MiddlewareFunc {
	c := DefaultConfig
	c.AccountID = accountID
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
	client := labstack.NewClient(config.AccountID, config.APIKey)
	cube := client.Cube()
	cube.APIKey = config.APIKey
	cube.BatchSize = config.BatchSize
	cube.DispatchInterval = config.DispatchInterval
	cube.ClientLookup = config.ClientLookup

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			// Start
			r := cube.Start(c.Request(), c.Response())

			// Handle panic
			defer func() {
				// Recover
				cube.Recover(recover(), r)

				// Stop
				cube.Stop(r, c.Response().Status, c.Response().Size)
			}()

			// Next
			if err = next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
