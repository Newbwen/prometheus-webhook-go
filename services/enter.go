package service

type ServiceGroup struct {
	DiskService
	K8sService
}

var ServiceGroupApp = new(ServiceGroup)
