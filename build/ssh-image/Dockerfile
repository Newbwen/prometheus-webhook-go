#job所使用的镜像
FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/alpine:latest
RUN  echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.20/main/" > /etc/apk/repositories \
&& echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.20/community/" >> /etc/apk/repositories \
&& apk update && apk upgrade \
&& apk add --no-cache openssh-client \
&& rm -rf /var/cache/apk/*
