# Kubernetes Node Cleaner Webhook

本项目是一个 **Prometheus Alertmanager Webhook 接收服务**，当收到节点磁盘使用率过高 (NodeDiskUsageHigh) 的告警时，会自动创建一个 **Kubernetes Job** 在指定节点执行清理任务。

## 功能特性

- 接收 Prometheus Alertmanager Webhook (`/webhook` 接口)。
- 从告警中提取节点 IP，自动查找对应的 NodeName。
- 在对应节点上创建清理 Job（通过 SSH 执行清理命令）。
- 成功的 Job 会在 **30 秒后自动删除**（避免资源堆积）。
- 失败的 Job 会保留，方便排查问题。

## 代码结构

```
.
├── main.go             # Webhook 主程序入口
├── go.mod              # Go 模块定义
├── go.sum
└── Dockerfile          # 构建镜像
```

## 依赖

- **Go**: 1.19+
- **Kubernetes client-go**
- **Gin Web 框架**


## 本地运行

### 1. 拉取依赖

```bash
go mod tidy
```

### 2. 运行服务

```bash
go run main.go
```

服务会监听 `:8080`，用于接收 Webhook：

### 3. 编译二进制文件
```bash
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"
go build -o prometheus-webhook main.go
```

### 4. 发送测试请求
- 发起磁盘使用率告警：
```bash
curl -XPOST http://localhost:8080/webhook -d '{"alerts":[{"status":"firing","labels":{"alertname":"NodeDiskUsageHigh","instance":"192.168.136.88:9100"},"annotations":{"summary":"Disk usage high"}}]}'
```
- 发起pod异常告警：
```bash
curl -XPOST http://localhost:8080/webhook -d '{
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "PodNotRunning",
        "namespace": "test",
        "pod": "nginx-6d98f5ffb4-kxbld"
      }
    },
    {
      "status": "firing",
      "labels": {
        "alertname": "PodNotRunning",
        "namespace": "test",
        "pod": "hotrod-5ccff444b7-54jbf"
      }
    }
  ]
}'
```

## Docker 打包

### 1. Dockerfile

```dockerfile
FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/busybox:latest

WORKDIR /app

COPY prometheus-webhook .

EXPOSE 8080

ENTRYPOINT ["./prometheus-webhook"]
```

### 2. 构建镜像

```bash
docker build -t my-registry/webhook-cleaner:latest .
```

### 3. 推送到镜像仓库

```bash
docker push my-registry/webhook-cleaner:latest
```

## 在 Kubernetes 部署

### 1. 部署 Webhook 服务

创建 `deployment.yaml`：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webhook-cleaner
  namespace: test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: webhook-cleaner
  template:
    metadata:
      labels:
        app: webhook-cleaner
    spec:
      serviceAccountName: default
      containers:
        - name: webhook-cleaner
          image: my-registry/webhook-cleaner:latest
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: prometheus-webhook
  namespace: test
spec:
  ports:
  - name: prometheus-webhook
    port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: prometheus-webhook
  type: NodePort
```

### 2. 配置 RBAC 权限
- 给予 `default` ServiceAccount 权限：给予 `job-creator` 角色，允许创建、获取、列出、监视、删除 Job。webhook创建job时使用
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: job-creator
rules:
- apiGroups: ["batch"]
  resources: ["jobs"]
  verbs: ["create", "get", "list", "watch", "delete"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: job-creator-binding
  namespace: default
subjects:
- kind: ServiceAccount
  name: default
  namespace: test
roleRef:
  kind: Role
  name: job-creator
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-reader
rules:
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["get", "list", "watch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: node-reader-binding
subjects:
- kind: ServiceAccount
  name: default
  namespace: test
roleRef:
  kind: ClusterRole
  name: node-reader
  apiGroup: rbac.authorization.k8s.io
```

- 给予 `default` ServiceAccount 权限：给予 `pod-manager` 角色，允许获取、列出、删除 Pod。webhook删除pod时使用
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: test
  name: pod-manager
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "delete"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: test
  name: pod-manager-binding
subjects:
- kind: ServiceAccount
  name: default
  namespace: test
roleRef:
  kind: Role
  name: pod-manager
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: pod-manager-global
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "delete"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: pod-manager-global-binding
subjects:
- kind: ServiceAccount
  name: default
  namespace: test
roleRef:
  kind: ClusterRole
  name: pod-manager-global
  apiGroup: rbac.authorization.k8s.io
```

### 3.部署configmap
```
#清理脚本
#!/bin/sh
NODE_IP=$1

ssh -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -i /root/.ssh/id_rsa root@$NODE_IP << 'EOF'
rm -rf /var/log/app/*

curl -X POST 192.168.136.88:6885/remsg -H "Content-Type: application/json" \
    --data '{"message":"应用日志清理完成!"}' \
EOF
#创建configmap
kubectl create configmap clean-script   --from-file=clean.sh
```

## 使用说明

- 当 Prometheus 触发 `NodeDiskUsageHigh`,`PodNotRunning` 告警时，Alertmanager 会调用该 Webhook。
- Webhook 服务会创建一个 Job，在目标节点通过 SSH 清理日志目录。收到pod异常告警时，Webhook会删除pod。
- 成功的 Job **30 秒后自动清理**，失败的 Job 会保留（方便排查）。
- prometheus-webhook会创建在test命名空间下。若是devops或者其他，则需要修改rbac权限
- shell脚本会post临时提交到NodePort，若是改为ingress需要修改其配置



