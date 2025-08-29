package service

type ServiceGroup struct {
	DiskService *DiskService
	K8sService  *K8sService
}

var ServiceGroupApp = new(ServiceGroup)

func Init() {
	//初始化k8s
	ServiceGroupApp.K8sService = NewK8sService()

}
