package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hongyuxuan/lizardrestic/common/errorx"
	"github.com/hongyuxuan/lizardrestic/server/internal/config"
	"github.com/hongyuxuan/lizardrestic/server/internal/handler"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	configFile       = kingpin.Flag("config", "config file").Short('f').Default("").String()
	etcdAddr         = kingpin.Flag("etcd-addr", "Etcd address.").Default("").String()
	servicePrefix    = kingpin.Flag("service-prefix", "Prefix of service key for registry. Can be empty").Default("").String()
	configurationDir = kingpin.Flag("configuration-dir", "Configuration directory.").Short('c').Default("").String()
	cacheDir         = kingpin.Flag("cache-dir", "Cache directory.").Default("").String()
	logLevel         = kingpin.Flag("log.level", "Log level.").Default("").String()
	listenOn         = kingpin.Flag("http-addr", "HTTP listen address.").Default("").String()
	metricsListenOn  = kingpin.Flag("metrics-addr", "Prometheus metrics listen address.").Default("").String()
	dbfile           = kingpin.Flag("db", "SQLite database file.").Short('d').Default("").String()

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
		etcdAddr,
		servicePrefix,
		configurationDir,
		cacheDir,
		listenOn,
		metricsListenOn,
		dbfile)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	ctx.SetVersion(AppVersion)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.LizardresticError:
			return e.Code, e.GetData()
		default:
			return http.StatusInternalServerError, errorx.HttpErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	})

	if c.Etcd.Address != "" {
		go handler.StartEtcdWatch(ctx)
	}

	// load cronjob
	restic := svc.NewResticService(context.Background(), ctx)
	restic.LoadCronJob()

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}

func printVersion() string {
	return fmt.Sprintf("App version: %s\nGo version:  %s\nBuild Time:  %s\nOS/Arch:     %s\nAuthor:      %s\n", AppVersion, GoVersion, BuildTime, OsArch, Author)
}
