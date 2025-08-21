package service

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sService struct{}

func (s *K8sService) CreateK8sJob(nodeIP, jobName string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	nodeName, err := s.findNodeNameByIP(clientset, nodeIP)
	if err != nil {
		return fmt.Errorf("查找 NodeName 失败: %v", err)
	}
	if nodeName == "" {
		return fmt.Errorf("未找到匹配的 NodeName for IP %s", nodeIP)
	}

	backoffLimit := int32(1)
	ttl := int32(30)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: jobName},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttl,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeName:      nodeName,
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{{
						Name:    "cleaner",
						Image:   "docker.io/library/ssh-client:v1.0", // TODO: 替换实际镜像
						Command: []string{"sh", "-c"},
						Args: []string{
							fmt.Sprintf(`ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i /root/.ssh/id_rsa root@%s "rm -rf /var/log/app/*"`, nodeIP),
						},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "ssh-key",
							MountPath: "/root/.ssh",
							ReadOnly:  true,
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "ssh-key",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/root/.ssh",
							},
						},
					}},
				},
			},
		},
	}
	_, err = clientset.BatchV1().Jobs("default").Create(context.Background(), job, metav1.CreateOptions{})
	return err
}

func (s *K8sService) findNodeNameByIP(clientset *kubernetes.Clientset, ip string) (string, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP && addr.Address == ip {
				return node.Name, nil
			}
		}
	}
	return "", nil
}
