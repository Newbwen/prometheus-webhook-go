package model

// Prometheus Alert 结构
type Alert struct {
	Status string `json:"status"`
	Labels struct {
		Alertname string `json:"alertname"`
		Instance  string `json:"instance"`
	} `json:"labels"`
	Annotations struct {
		Summary string `json:"summary"`
	} `json:"annotations"`
}

type WebhookPayload struct {
	Alerts []Alert `json:"alerts"`
}
