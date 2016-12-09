package plugin

import "github.com/labstack/echo"

// Add/remove HTTP response headers.

type (
	Header struct {
		Base `json:",squash"`
		Set  map[string]string `json:"set"`
		Add  map[string]string `json:"add"`
		Del  []string          `json:"del"`
	}
)

func (*Header) Init() (err error) {
	return
}

func (*Header) Priority() int {
	return 1
}

func (h *Header) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Response().Header()
		for k, v := range h.Set { // Set headers
			header.Set(k, v)
		}
		for k, v := range h.Add { // Add headers
			header.Add(k, v)
		}
		for _, k := range h.Del { // Delete headers
			header.Del(k)
		}
		return next(c)
	}
}
