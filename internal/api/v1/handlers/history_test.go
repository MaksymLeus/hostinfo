package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hostinfo/internal/metrics"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHistoryHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call History handler
	if assert.NoError(t, History(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var points []metrics.DataPoint
		err := json.Unmarshal(rec.Body.Bytes(), &points)
		assert.NoError(t, err)

		// The history buffer records at least 1 point on startup/collection
		assert.NotNil(t, points)
	}
}
