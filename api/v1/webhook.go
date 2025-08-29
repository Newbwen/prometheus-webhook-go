package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Newbwen/prometheus-webhook-go/model"
	"github.com/Newbwen/prometheus-webhook-go/service"
	"github.com/gin-gonic/gin"
)

type WebhookApi struct{}

func (a *WebhookApi) HealthzCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// 接收返回的信息
func (a *WebhookApi) ReMessage(c *gin.Context) {
	var reMessage map[string]interface{}
	if err := c.ShouldBindJSON(&reMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("脚本执行成功，删除文件...", reMessage["message"])
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

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
		if alert.Status == "firing" && alert.Labels["alertname"] == "NodeDiskUsageHigh" {
			go service.ServiceGroupApp.DiskService.HandleDiskAlert(alert.Labels["instance"])
		}
		if alert.Status == "firing" && alert.Labels["alertname"] == "PodNotRunning" {
			pod := alert.Labels["pod"]
			namespace := alert.Labels["namespace"]
			log.Printf("收到 PodNotRunning 告警, 重启 Pod %s/%s\n", namespace, pod)
			go service.ServiceGroupApp.K8sService.DeletePod(namespace, pod)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
