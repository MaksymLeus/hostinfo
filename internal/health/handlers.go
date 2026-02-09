package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Health(c echo.Context) error {
	// return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	return c.JSON(http.StatusOK, CheckResult{
		Status: "ok",
	})
}

func Live(c echo.Context) error {
	checks := BasicChecks()

	return c.JSON(http.StatusOK, CheckResult{
		Status: "ok",
		Checks: checks,
	})
}

func Ready(c echo.Context) error {
	return c.JSON(http.StatusOK, CheckResult{
		Status: "ok",
	})
}
