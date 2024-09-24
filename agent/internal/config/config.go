package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	ServicePrefix    string `json:","`
	ConfigurationDir string `json:","`
	CacheDir         string `json:","`
}

func NewConfig(configFile, logLevel, etcdHost, serviceKey, servicePrefix, configurationDir, cacheDir, listenOn, metricsListenOn *string) Config {
	var c = Config{}
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		c.Name = "LizardAgent"
		c.ListenOn = "0.0.0.0:5017"
		c.Timeout = 60000
		c.Log.Encoding = "plain"
		c.Log.Level = "info"
		c.Prometheus.Host = "0.0.0.0"
		c.Prometheus.Port = 15017
		c.Prometheus.Path = "/metrics"
		c.Etcd = discov.EtcdConf{}
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *servicePrefix != "" {
		c.ServicePrefix = *servicePrefix
	}
	if *etcdHost != "" {
		c.Etcd.Hosts = strings.Split(*etcdHost, ",")
		if *serviceKey != "" {
			c.Etcd.Key = *serviceKey
		}
	}
	if *configurationDir != "" {
		c.ConfigurationDir = *configurationDir
	}
	if *cacheDir != "" {
		c.CacheDir = *cacheDir
	}
	c.Etcd.Key = c.ServicePrefix + c.Etcd.Key
	if *listenOn != "" {
		c.ListenOn = *listenOn
	}
	if *metricsListenOn != "" {
		arr := strings.Split(*metricsListenOn, ":")
		c.Prometheus.Host = arr[0]
		port, _ := strconv.Atoi(arr[1])
		c.Prometheus.Port = port
	}
	if len(c.Etcd.Hosts) == 0 {
		logx.Errorf("Etcd host must be specified.")
		os.Exit(0)
	}
	logx.DisableStat()
	logx.MustSetup(c.Log)
	logx.Infof("Using config: %+v", c)
	return c
}
