package restic

import (
	"context"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LssnapshotsLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	commonRestic *commonsvc.CommonRestic
}

func NewLssnapshotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LssnapshotsLogic {
	return &LssnapshotsLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *LssnapshotsLogic) Lssnapshots(req *types.LsSnapshotsReq) (resp *types.Response, err error) {
	var output string
	var args = []string{"ls", req.SnapshotId, req.Dir, "--host", req.Host}
	if req.Tag != "" {
		args = append(args, "--tag", req.Tag)
	}
	if output, err = l.commonRestic.RunCommand(req.RepoUrl, args...); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: output,
	}
	return
}
