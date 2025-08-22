package model

// Prometheus Alert 结构
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"` // 开始时间
	EndsAt       string            `json:"endsAt"`   // 结束时间
	GeneratorURL string            `json:"generatorURL"`
}

type WebhookPayload struct {
	Alerts []Alert `json:"alerts"`
}
