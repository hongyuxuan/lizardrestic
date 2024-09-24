package restic

import (
	"context"
	"encoding/json"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsrepoLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	commonRestic *commonsvc.CommonRestic
}

func NewStatsrepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsrepoLogic {
	return &StatsrepoLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *StatsrepoLogic) Statsrepo(req *types.StatsRepoReq) (resp *types.Response, err error) {
	var output string
	var args = []string{"stats", "--json"}
	if output, err = l.commonRestic.RunCommand(req.RepoUrl, args...); err != nil {
		return
	}
	var stats commontypes.RepositoryStats
	if err = json.Unmarshal([]byte(output), &stats); err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: stats,
	}
	return
}
