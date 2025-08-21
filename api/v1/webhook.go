package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Newbwen/prometheus-webhook-go/model"
	"github.com/Newbwen/prometheus-webhook-go/services"
	"github.com/gin-gonic/gin"
)

type WebhookApi struct{}

func (a *WebhookApi) HandleWebhook(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("read body error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	var payload model.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("json unmarshal error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	for _, alert := range payload.Alerts {
		if alert.Status == "firing" && alert.Labels.Alertname == "NodeDiskUsageHigh" {
			go service.ServiceGroupApp.DiskService.HandleDiskAlert(alert.Labels.Instance)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
