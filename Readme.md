# 服务监控组件

## 目录结构

- 依赖组件:
    * `docker` 
    * `docker-compose` 

- 需要主机开启对外端口：
    * 9090 (`prometheus`)
    * 3000 (`grafana`)

- 全服统一配置文件:
    * `config/scene/StartSceneConfig.txt`
    * `config/scene/StartProcessConfig.txt`
    * `config/scene/StartMachineConfig.txt`

- 执行脚本开启后可访问`http://localhost:3000`查看服务状态

- prometheus服务配置文件:
    * `config/prometheus/prometheus.yml`

## 开启监控服务
1. 更改`StartSceneConfig.txt`，`StartProcessConfig.txt`，`StartMachineConig.txt`后，更新到线上服务器的`config/scene`目录中
2. 执行不同系统对应的`start`脚本，脚本会将`txt`配置文件合并转换为`prometheus`可读取的配置文件`prometheus.yml`，并且开启`prometheus`和`grafana`的`docker`镜像，此时访问`http://localhost:3000`即可使用监控功能
3. 配置文件有更改后，直接运行`start`脚本，监控服务将重启，不会对其他线上服务造成任何影响
> 注意：在本机测试时如果服务的`InnerIP`配置为`127.0.0.1`，则需要将`config/consul/service.json`文件中所有`127.0.0.1`的ip地址修改为本地局域网地址例如`192.168.1.137`，因为在`docker`镜像中解析`localhost`地址解析不到宿主机中，所以本地测试时最好将`InnerIP`配置为本机局域网ip地址