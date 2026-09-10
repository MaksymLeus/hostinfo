package handlers

import (
	"hostinfo/internal/metrics"
	"net/http"

	"github.com/labstack/echo/v4"
)

// History godoc
// @Summary Get historical metrics
// @Description Returns recent in-memory CPU and memory time-series data for sparklines
// @Tags info
// @Accept json
// @Produce json
// @Success 200 {array} metrics.DataPoint
// @Router /history [get]
func History(c echo.Context) error {
	points := metrics.History.GetHistory()
	return c.JSON(http.StatusOK, points)
}
