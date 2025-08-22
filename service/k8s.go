package service

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

type K8sService struct {
	Clientset *kubernetes.Clientset
}

// 初始化k8s客户端
func NewK8sService() *K8sService {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("failed to get in-cluster config: %v", err))
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create k8s clientset: %v", err))
	}
	return &K8sService{Clientset: clientset}
}

func (s *K8sService) CreateK8sJob(nodeIP, jobName string) error {

	nodeName, err := s.findNodeNameByIP(s.Clientset, nodeIP)
	if err != nil {
		return fmt.Errorf("查找 NodeName 失败: %v", err)
	}
	if nodeName == "" {
		return fmt.Errorf("未找到匹配的 NodeName for IP %s", nodeIP)
	}

	backoffLimit := int32(1)
	ttl := int32(30)
	// 创建 Job
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
						Image:   "docker.io/library/ssh-client:v1.0",
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
	_, err = s.Clientset.BatchV1().Jobs("default").Create(context.Background(), job, metav1.CreateOptions{})
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
func (s *K8sService) DeletePod(namespace, name string) error {
	err := s.Clientset.CoreV1().Pods(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("删除 Pod %s 失败: %v\n", name, err)
		return err
	}
	fmt.Printf("删除 Pod %s 成功\n", name)
	return nil
}
