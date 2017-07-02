package labstack

import (
	"time"

	"github.com/dghubble/sling"
)

type (
	Client struct {
		sling *sling.Sling
		// LoggingBatchSize        int
		// LoggingDispatchInterval int
	}

	Error struct {
	}

	Log struct {
		ID      string    `json:"id,omitempty"`
		Time    time.Time `json:"time"`
		Module  string    `json:"module"`
		Level   string    `json:"level"`
		Message string    `json:"message"`
	}
)

const (
	apiURL = "https://api.labstack.com"
)

func NewClient(apiKey string) *Client {
	return &Client{
		sling: sling.New().Base(apiURL).Add("Authorization", "Bearer "+apiKey),
	}
}

func (c *Client) WriteLog(l *Log) (err error) {
	_, err = c.sling.New().Post("/logging").BodyJSON(l).ReceiveSuccess(l)
	return
}
