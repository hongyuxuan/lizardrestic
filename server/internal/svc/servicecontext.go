package svc

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardrestic/agent/lizardagent"
	"github.com/hongyuxuan/lizardrestic/common/errorx"
	"github.com/hongyuxuan/lizardrestic/common/utils"
	"github.com/hongyuxuan/lizardrestic/server/internal/config"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	AgentList  map[string]*types.RpcAgent
	EtcdClient *clientv3.Client
	Sqlite     *gorm.DB
	Version    string
	Cron       *cron.Cron
	CronIdMap  map[string]cron.EntryID
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config:    c,
		AgentList: make(map[string]*types.RpcAgent),
		Sqlite:    utils.NewSQLite(c.Sqlite, c.Log.Level),
		Cron:      cron.New(cron.WithSeconds()),
		CronIdMap: make(map[string]cron.EntryID),
	}
	if c.Etcd.Address != "" {
		etcdHosts := strings.Split(c.Etcd.Address, ",")
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   etcdHosts,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			logx.Errorf("Failed to connect to etcd: %v", err)
			os.Exit(0)
		}
		logx.Infof("Connect to etcd host=%s success", c.Etcd.Address)
		svcCtx.EtcdClient = client
	}

	return svcCtx
}

func (s *ServiceContext) GetTargetAgent(ip string) (agent lizardagent.LizardAgent, err error) {
	for k, v := range s.AgentList {
		re, _ := regexp.Compile(fmt.Sprintf("%slizardrestic-agent\\.(.+?)\\.%s", s.Config.ServicePrefix, ip))
		if re.MatchString(k) {
			return v.Client, nil
		}
	}
	return nil, errorx.NewDefaultError("Cannot find lizardrestic-agent of ip=%s, maybe the server cannot communicated with the agent", ip)
}

func (s *ServiceContext) SetVersion(version string) {
	s.Version = version
}
