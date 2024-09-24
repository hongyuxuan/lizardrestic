package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Etcd             EtcdConf `json:",optional"`
	ServicePrefix    string   `json:",optional"`
	Sqlite           string
	Rpc              RpcOption `json:",optional"`
	ConfigurationDir string    `json:","`
	CacheDir         string    `json:","`
}

type RpcOption struct {
	Timeout       int64 `json:",optional"`
	KeepaliveTime int64 `json:",optional"`
	RetryInterval int64 `json:",optional"`
}

func (rpc RpcOption) IsEmpty() bool {
	return reflect.DeepEqual(rpc, RpcOption{})
}

type EtcdConf struct {
	Address string
}

func NewConfig(configFile, logLevel, etcdAddr, servicePrefix, configurationDir, cacheDir, listenOn, metricsListenOn, dbfile *string) Config {
	var c = Config{}
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		c.Name = "LizardServer"
		c.Host = "0.0.0.0"
		c.Port = 5117
		c.Timeout = 60000
		c.Log.Encoding = "plain"
		c.Log.Level = "info"
		c.Prometheus.Host = "0.0.0.0"
		c.Prometheus.Port = 15117
		c.Prometheus.Path = "/metrics"
		c.Etcd = EtcdConf{}
		c.Sqlite = "./lizardcd.db"
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *etcdAddr != "" {
		c.Etcd.Address = *etcdAddr
	}
	if *configurationDir != "" {
		c.ConfigurationDir = *configurationDir
	}
	if *cacheDir != "" {
		c.CacheDir = *cacheDir
	}
	if *listenOn != "" {
		arr := strings.Split(*listenOn, ":")
		c.Host = arr[0]
		port, _ := strconv.Atoi(arr[1])
		c.Port = port
	}
	if *metricsListenOn != "" {
		arr := strings.Split(*metricsListenOn, ":")
		c.Prometheus.Host = arr[0]
		port, _ := strconv.Atoi(arr[1])
		c.Prometheus.Port = port
	}
	if *servicePrefix != "" {
		c.ServicePrefix = *servicePrefix
	}
	if *dbfile != "" {
		c.Sqlite = *dbfile
	}
	if c.Rpc.IsEmpty() {
		c.Rpc = RpcOption{}
	}
	if c.Rpc.KeepaliveTime == 0 {
		c.Rpc.KeepaliveTime = 600
	}
	if c.Rpc.Timeout == 0 {
		c.Rpc.Timeout = 2000
	}
	if c.Rpc.RetryInterval == 0 {
		c.Rpc.RetryInterval = 10
	}
	if c.Etcd.Address == "" {
		logx.Errorf("Etcd address must be specified.")
		os.Exit(0)
	}
	logx.DisableStat()
	logx.MustSetup(c.Log)
	logx.Infof("Using config: %+v", c)
	return c
}
