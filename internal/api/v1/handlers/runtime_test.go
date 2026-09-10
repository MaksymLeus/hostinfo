package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDetectCloudHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cloud", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Cloud(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var info CloudInfo
		err := json.Unmarshal(rec.Body.Bytes(), &info)
		assert.NoError(t, err)
		// When running locally outside AWS/GCP/Azure, it defaults to provider "local"
		assert.NotEmpty(t, info.Provider)
	}
}

func TestDetectRuntimeHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/kubernetes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Kubernetes(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var info KubernetesInfo
		err := json.Unmarshal(rec.Body.Bytes(), &info)
		assert.NoError(t, err)
		// Outside K8s, enabled should be false by default
		assert.False(t, info.Enabled)
	}
}
