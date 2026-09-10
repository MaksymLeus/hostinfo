package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPingMissingParam(t *testing.T) {
	e := echo.New()
	// Request without "host" query parameter
	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Should fail with 400 Bad Request
	if assert.NoError(t, Ping(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "host")
	}
}

func TestCurlMissingParam(t *testing.T) {
	e := echo.New()
	// Request without "url" query parameter
	req := httptest.NewRequest(http.MethodPost, "/curl", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Should fail with 400 Bad Request
	if assert.NoError(t, Curl(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "url")
	}
}

func TestDigMissingParam(t *testing.T) {
	e := echo.New()
	// Request without "host" query parameter
	req := httptest.NewRequest(http.MethodPost, "/dns", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Should fail with 400 Bad Request
	if assert.NoError(t, Dig(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "host")
	}
}

func TestTCPMissingParam(t *testing.T) {
	e := echo.New()
	// Request without host or port
	req := httptest.NewRequest(http.MethodPost, "/tcp", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Should fail with 400 Bad Request
	if assert.NoError(t, TCP(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "host")
	}
}
