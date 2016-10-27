package plugin

// import (
// 	"io/ioutil"
// 	"net/http"
//
// 	"github.com/golang/protobuf/proto"
// 	"github.com/labstack/echo"
// 	"github.com/nats-io/nats"
// )
//
// type (
// 	NATS struct {
// 		Base    `json:",squash"`
// 		Subject string `json:"subject"`
// 		// Publish *natsPub `json:"publish"`
// 		// Request *natsReq `json:"request"`
// 		Async bool `json:"async"`
// 		conn  *nats.Conn
// 	}
//
// 	// natsPub struct {
// 	// 	Subject string `json:"subject"`
// 	// }
// 	//
// 	// natsReq struct {
// 	// 	Timeout time.Duration `json:"timeout"`
// 	// }
// )
//
// func (n *NATS) Initialize() (err error) {
// 	// Defaults
// 	// if n.Request != nil {
// 	// 	// TODO: https://github.com/nats-io/nats/blob/master/TODO.md
// 	// 	if n.Request.Timeout == 0 {
// 	// 		n.Request.Timeout = 1 * time.Minute
// 	// 	}
// 	// }
//
// 	// Initialize
// 	if n.conn, err = nats.Connect(nats.DefaultURL); err != nil {
// 		return
// 	}
// 	// TODO: defer n.conn.Close()
// 	return
// }
//
// func (*NATS) Priority() int {
// 	return 1
// }
//
// func (n *NATS) Process(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) (err error) {
// 		// Message
// 		body := c.Request().Body
// 		msg := new(Message)
// 		msg.ExternalId = "xxx"
// 		msg.Body, err = ioutil.ReadAll(body) // Body
// 		if err != nil {
// 			return
// 		}
// 		b, err := proto.Marshal(msg)
// 		if err != nil {
// 			return err
// 		}
//
// 		// Publish
// 		// TODO: in a goroutine?
// 		// if n.Publish != nil {
// 		if err = n.conn.Publish(n.Subject, b); err != nil {
// 			return
// 		}
// 		// }
//
// 		// // Request/reply
// 		// if n.Request != nil {
// 		// 	rep, err := n.conn.Request(n.Request.Subject, b, n.Request.Timeout)
// 		// 	if err != nil {
// 		// 		return err
// 		// 	}
// 		// 	msg := new(Message)
// 		// 	if err := proto.Unmarshal(rep.Data, msg); err != nil {
// 		// 		return err
// 		// 	}
// 		// 	for k, v := range msg.Header { // Copy headers
// 		// 		c.Response().Header().Add(k, v)
// 		// 	}
// 		// 	c.Response().WriteHeader(200)
// 		// 	_, err = c.Response().Write(msg.Body)
// 		// 	return err
// 		// }
//
// 		return c.NoContent(http.StatusOK)
// 	}
// }
