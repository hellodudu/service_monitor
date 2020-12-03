#!/bin/bash
echo "停止prometheus服务..."
docker-compose down

echo "转换txt配置文件为prometheus.yml..."
./bin/exporter_mac -import_path=./config/scene/ -export_path=./config/prometheus/

echo "启动prometheus容器..."
docker-compose up -d

echo "success..."