# Lizardrestic - 基于 Restic 的企业级文件备份平台
[restic](https://github.com/restic/restic) 是一个单机版的文件备份工具，支持将本地文件备份到 S3 上。

`restic` 只是一个命令行程序，不提供多服务器的远程备份管理，因此它只是一个工具，而不是一个平台。

Lizardrestic 在 `restic` 的基础上封装了一层，提供了远程管理、定时任务等功能，形成了一个真正企业可用的备份平台。

## 功能特性
- agent 管理：每个服务器上需安装 agent（已内置restic-0.17.0），并自动注册。注册功能依赖 `etcd`。
- 仓库管理：lizardrestic 支持将数据备份到 S3，可在管理界面添加多个 S3 仓库。备份时支持选择仓库。
- 策略管理：可创建备份策略，可选择仓库、设置定时任务、设置标签、指定备份目录、排除路径、备份快照保留时间、选择备份主机、是否启用。
- 快照管理：查看快照、从快照中搜索文件、列出文件、手动删除快照、指定快照恢复到某路径。
- 执行记录：备份/恢复任务历史记录，任务执行结果查看等。

## 部署
### migrate
在 [release](https://github.com/hongyuxuan/lizardrestic/releases) 页面下载 migrate 安装包，解压：
```shell
tar zxf migrate-linux-amd64-<version>.tar.gz
```

执行 `migrate` 命令初始化数据库、创建相关表。
```shell
./migrate -d /path/to/your/lizardrestic.db
```
注意 `lizardrestic.db` 必须和 server 启动用的是同一个 db 文件。

### server
在 [release](https://github.com/hongyuxuan/lizardrestic/releases) 页面下载 server 安装包，解压：
```shell
tar zxf lizardrestic-server-linux-amd64-<version>.tar.gz
cd lizardrestic-server-linux-amd64-<version> && tree .
.
├── bin
│   ├── lizardrestic-server
│   └── restic
├── cache
├── configuration
├── docs
│   └── swagger.json
└── etc
    └── resticserver.yaml
```
`bin` 下面有两个可执行文件：`lizardrestic-server` 是主程序，`restic` 是备份命令行工具。

编辑配置文件 etc/resticserver.yaml
```yaml
Name: lizardServer
Host: 0.0.0.0
Port: 7138
Timeout: 60000
Log:
  Encoding: plain
  Level: debug
Prometheus:
  Host: 0.0.0.0
  Port: 17138
  Path: /metrics
Etcd:
  Address: 10.50.89.17:2379
Sqlite: ./lizardrestic.db
ServicePrefix: backup-
Rpc:
  Timeout: 2000 # millisecond
  KeepaliveTime: 600 # seconds
  RetryInterval: 600 # seconds
ConfigurationDir: ./configuration
CacheDir: ./cache
```
`ConfigurationDir` 存放 server 连接 S3 的配置文件。请改成实际目录。

`CacheDir` 存放备份文件的本地缓存。请改成实际目录。

另外，需要提前准好一个 `etcd`。lizardrestic 依赖 `etcd` 做 agent 自动发现。

启动 lizardrestic-server
```shell
./bin/lizardrestic-server -f etc/resticserver.yaml -d lizardrestic.db
```



### ui
目前 ui 仅提供 docker 镜像。用源码里的 [deployment/manifests/nginx.conf](./deployment/manifests/nginx.conf) 替换镜像里的 `/etc/nginx.conf`，并修改里面的 server 地址为实际地址。

### agent
在 [release](https://github.com/hongyuxuan/lizardrestic/releases) 页面下载 agent 安装包，解压：
```shell
tar zxf lizardrestic-agent-linux-amd64-<version>.tar.gz
cd lizardrestic-agent-linux-amd64-<version> && tree .
.
├── bin
│   ├── lizardrestic-agent
│   └── restic
├── cache
├── configuration
└── etc
    └── agent.yaml
```

编辑配置文件 etc/agent.yaml
```yaml
Name: lizardAgent
ListenOn: 0.0.0.0:7038
Timeout: 60000
Log:
  Encoding: plain
  Level: info 
Prometheus:
  Host: 0.0.0.0
  Port: 17038
  Path: /metrics
Etcd:
  Hosts:
  - 10.50.89.17:2379
  Key: lizardrestic-agent.default.127.0.0.1
ServicePrefix: backup-
ConfigurationDir: ./configuration
CacheDir: ./cache
```
将 `etcd` 注册 key 里 `127.0.0.1` 改成实际的 agent IP 地址。

保持 `ServicePrefix` 和 server 的配置文件一样。

启动 lizardrestic-agent
```shell
./bin/lizardrestic-agent -f etc/agent.yaml
```

打开 ui 页面，切换到 `Agent列表`，观察 agent 是否已经自动发现。