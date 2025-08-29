#!/bin/sh
NODE_IP=$1

ssh -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -i /root/.ssh/id_rsa root@$NODE_IP << 'EOF'
# 清理应用日志
LOGS=`ls /var/log/app`
rm -rf /var/log/app/*
echo "======目录清理完成======"

curl -X POST prometheus-webhook.svc.cluster.local -H "Content-Type: application/json" \
    --data '{"message":"应用日志清理完成\n$LOGS"}' \
EOF