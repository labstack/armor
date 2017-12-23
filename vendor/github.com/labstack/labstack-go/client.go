package labstack

import (
	"sync"

	"github.com/dghubble/sling"
	glog "github.com/labstack/gommon/log"
)

type (
	Client struct {
		accountID string
		apiKey    string
		sling     *sling.Sling
		logger    *glog.Logger
	}

	Fields map[string]interface{}

	SearchParameters struct {
		Query       string   `json:"query"`
		QueryString string   `json:"query_string"`
		Since       string   `json:"since"`
		Sort        []string `json:"sort"`
		Size        int      `json:"size"`
		From        int      `json:"from"`
	}
)

const (
	apiURL = "https://api.labstack.com"
)

// NewClient creates a new client for the LabStack API.
func NewClient(accountID, apiKey string) *Client {
	return &Client{
		accountID: accountID,
		apiKey:    apiKey,
		sling:     sling.New().Base(apiURL).Add("Authorization", "Bearer "+apiKey),
		logger:    glog.New("labstack"),
	}
}

// Cube returns the cube service.
func (c *Client) Cube() (cube *Cube) {
	cube = &Cube{
		sling:            c.sling.Path("/cube"),
		mutex:            new(sync.RWMutex),
		logger:           c.logger,
		AccountID:        c.accountID,
		APIKey:           c.apiKey,
		BatchSize:        60,
		DispatchInterval: 60,
	}
	cube.resetRequests()
	return
}

// Jet returns the jet service.
func (c *Client) Jet() *Jet {
	return &Jet{
		sling:  c.sling.Path("/jet"),
		logger: c.logger,
	}
}
