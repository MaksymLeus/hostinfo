package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health godoc
// @Summary Basic health check
// @Description Returns basic service status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} CheckResult
// @Router /healthz [get]
func Health(c echo.Context) error {
	// return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	return c.JSON(http.StatusOK, CheckResult{
		Status: "ok",
	})
}

// Live godoc
// @Summary Liveness probe
// @Description Returns service liveness status and uptime check
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} CheckResult
// @Router /healthz/live [get]
func Live(c echo.Context) error {
	checks := BasicChecks()

	return c.JSON(http.StatusOK, CheckResult{
		Status: "ok",
		Checks: checks,
	})
}

// Ready godoc
// @Summary Readiness probe
// @Description Returns service readiness status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} CheckResult
// @Router /healthz/ready [get]
func Ready(c echo.Context) error {
	return c.JSON(http.StatusOK, CheckResult{
		Status: "ok",
	})
}
