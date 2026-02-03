package health

type CheckResult struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks,omitempty"`
}
