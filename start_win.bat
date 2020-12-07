@echo off
REM 后续命令使用的是：UTF-8编码
chcp 65001
echo .
echo 停止prometheus服务...
docker-compose down

REM "转换txt配置文件为prometheus.yml..."
.\bin\exporter_win.exe -import_path=.\config\scene\ -export_path=.\config\prometheus\

echo .
echo 启动prometheus容器...
docker-compose up -d

echo .
echo success...
pause