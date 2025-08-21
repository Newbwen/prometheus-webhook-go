package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Newbwen/prometheus-webhook-go/services"
	"github.com/gin-gonic/gin"
)

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

type AlertApi struct{}

// Webhook 处理函数
func (a *AlertApi) WebhookHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("read body error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("json unmarshal error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	for _, alert := range payload.Alerts {
		if alert.Status == "firing" && alert.Labels.Alertname == "NodeDiskUsageHigh" {
			go service.HandleDiskAlert(alert.Labels.Instance)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
