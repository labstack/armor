package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := func(c echo.Context) error {
		sess, _ := Get("test", c)
		sess.Options.Domain = "labstack.com"
		sess.Values["foo"] = "bar"
		sess.Save(c.Request(), c.Response())
		return c.String(http.StatusOK, "test")
	}
	store := sessions.NewCookieStore([]byte("secret"))
	config := Config{
		Skipper: func(c echo.Context) bool {
			return true
		},
		Store: store,
	}

	// Skipper
	mw := MiddlewareWithConfig(config)
	h := mw(echo.NotFoundHandler)
	h(c)
	assert.Nil(t, c.Get(key))

	// Panic
	config.Skipper = nil
	config.Store = nil
	assert.Panics(t, func() {
		MiddlewareWithConfig(config)
	})

	// Core
	mw = Middleware(store)
	h = mw(handler)
	h(c)
	assert.Contains(t, rec.Header().Get(echo.HeaderSetCookie), "labstack.com")
}
