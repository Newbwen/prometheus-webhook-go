package service

import (
	"fmt"
	"strings"
)

type DiskService struct{}

func (s *DiskService) HandleDiskAlert(instance string) {
	nodeIP := strings.Split(instance, ":")[0]
	jobName := fmt.Sprintf("clean-node-%s", strings.ReplaceAll(nodeIP, ".", "-"))

	fmt.Println("创建清理 Job:", jobName, "节点:", nodeIP)

	if err := ServiceGroupApp.K8sService.CreateK8sJob(nodeIP, jobName); err != nil {
		fmt.Println("创建 Job 失败:", err)
	} else {
		fmt.Println("Job 创建成功:", jobName)
	}
}
