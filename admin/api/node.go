package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Node struct {
		Name string `json:"name"`
	}
)

func (h *handler) nodes(c echo.Context) error {
	cluster := h.armor.Cluster
	nodes := []*Node{}

	for _, m := range cluster.Members() {
		nodes = append(nodes, &Node{
			Name: m.Name,
		})
	}

	return c.JSON(http.StatusOK, nodes)
}
