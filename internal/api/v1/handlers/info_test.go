package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInfoHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call Info handler
	if assert.NoError(t, Info(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var info HostInfo
		err := json.Unmarshal(rec.Body.Bytes(), &info)
		assert.NoError(t, err)

		// Validate that core system info is populated
		assert.NotEmpty(t, info.OS)
		assert.NotEmpty(t, info.Arch)
		assert.NotEmpty(t, info.GoVersion)
		assert.NotEmpty(t, info.Uptime)

		// Validate hardware stats are populated
		assert.NotNil(t, info.CPU)
		assert.NotNil(t, info.Memory)
	}
}
