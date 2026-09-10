package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call Health handler
	if assert.NoError(t, Health(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var res CheckResult
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "ok", res.Status)
	}
}

func TestLiveHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call Live handler
	if assert.NoError(t, Live(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var res CheckResult
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "ok", res.Status)
		assert.NotEmpty(t, res.Checks["uptime"])
	}
}
