package plugin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestStatic(t *testing.T) {
	e := echo.New()
	s := &Static{
		Root: "../_fixture",
	}
	s.Init()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// File found
	c.SetParamNames("*")
	c.SetParamValues("/images/walle.png")
	h := s.Process(echo.NotFoundHandler)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Header().Get(echo.HeaderContentLength), "219885")
	}

	// File not found
	rec.Body.Reset()
	c.SetParamNames("*")
	c.SetParamValues("none")
	h = s.Process(echo.NotFoundHandler)
	he := h(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusNotFound, he.Code)

	// HTML5
	rec.Body.Reset()
	s.HTML5 = true
	c.SetParamNames("*")
	c.SetParamValues("random")
	h = s.Process(echo.NotFoundHandler)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Armor")
	}

	// Browse
	rec.Body.Reset()
	s.Browse = true
	c.SetParamNames("*")
	c.SetParamValues("")
	h = s.Process(echo.NotFoundHandler)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "images")
	}
}
