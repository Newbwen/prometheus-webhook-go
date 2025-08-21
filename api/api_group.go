package api

type PrometheusApiGroup struct {
	AlertApi AlertApi
}

type ApiGroup struct {
	PrometheusApiGroup PrometheusApiGroup
}

var ApiGroupApp = new(ApiGroup)
