FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/busybox:latest

WORKDIR /app

COPY prometheus-webhook .

EXPOSE 8080

ENTRYPOINT ["./prometheus-webhook"]
