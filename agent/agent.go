package main

import (
	"fmt"

	"github.com/hongyuxuan/lizardrestic/agent/internal/config"
	"github.com/hongyuxuan/lizardrestic/agent/internal/server"
	"github.com/hongyuxuan/lizardrestic/agent/internal/svc"
	"github.com/hongyuxuan/lizardrestic/agent/types/agent"

	"github.com/alecthomas/kingpin/v2"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configFile       = kingpin.Flag("config", "config file").Short('f').Default("").String()
	logLevel         = kingpin.Flag("log.level", "Log level.").Default("").String()
	etcdHost         = kingpin.Flag("etcd-host", "Etcd hosts.").Default("").String()
	serviceKey       = kingpin.Flag("service-key", "Service key for registry. Format must be: lizardrestic-agent.<system>.<ip>").Default("").String()
	servicePrefix    = kingpin.Flag("service-prefix", "Prefix of service key for registry. Can be empty").Default("").String()
	configurationDir = kingpin.Flag("configuration-dir", "Configuration directory.").Short('c').Default("").String()
	cacheDir         = kingpin.Flag("cache-dir", "Cache directory.").Default("").String()
	listenOn         = kingpin.Flag("grpc-addr", "Grpc listen address.").Default("").String()
	metricsListenOn  = kingpin.Flag("metrics-addr", "Prometheus metrics listen address.").Default("").String()

	/* print app version */
	AppVersion = "unknown"
	GoVersion  = "unknown"
	BuildTime  = "unknown"
	OsArch     = "unknown"
	Author     = "unknown"
)

func main() {
	// Parse flags
	kingpin.Version(printVersion())
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	c := config.NewConfig(
		configFile,
		logLevel,
		etcdHost,
		serviceKey,
		servicePrefix,
		configurationDir,
		cacheDir,
		listenOn,
		metricsListenOn)

	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		agent.RegisterLizardAgentServer(grpcServer, server.NewLizardAgentServer(ctx))
		logx.Infof("Lizardrestic-agent: %s register to etcd success", c.Etcd.Key)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}

func printVersion() string {
	return fmt.Sprintf("App version: %s\nGo version:  %s\nBuild Time:  %s\nOS/Arch:     %s\nAuthor:      %s\n", AppVersion, GoVersion, BuildTime, OsArch, Author)
}
