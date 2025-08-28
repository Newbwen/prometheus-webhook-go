package service

import (
	"context"
	//batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

// Test findNodeNameByIP 能找到 Node
func TestFindNodeNameByIP(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node-1"},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "192.168.1.10"},
			},
		},
	})

	svc := &K8sService{Clientset: fakeClient}
	nodeName, err := svc.findNodeNameByIP(fakeClient, "192.168.1.10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nodeName != "node-1" {
		t.Fatalf("expected node-1, got %s", nodeName)
	}
}

// Test CreateK8sJob 能成功创建 Job
func TestCreateK8sJob(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node-2"},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "192.168.1.20"},
			},
		},
	})

	svc := &K8sService{Clientset: fakeClient}

	err := svc.CreateK8sJob("192.168.1.20", "test-job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证 Job 是否真的创建成功
	jobs, err := fakeClient.BatchV1().Jobs("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error listing jobs: %v", err)
	}
	if len(jobs.Items) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Items))
	}
	if jobs.Items[0].Name != "test-job" {
		t.Fatalf("expected job name test-job, got %s", jobs.Items[0].Name)
	}
}

// Test DeletePod 能删除 Pod
func TestDeletePod(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	})

	svc := &K8sService{Clientset: fakeClient}
	err := svc.DeletePod("default", "test-pod")
	if err != nil {
		t.Fatalf("unexpected error deleting pod: %v", err)
	}

	// 确认 Pod 被删除
	_, err = fakeClient.CoreV1().Pods("default").Get(context.Background(), "test-pod", metav1.GetOptions{})
	if err == nil {
		t.Fatalf("expected pod to be deleted, but still exists")
	}
}
