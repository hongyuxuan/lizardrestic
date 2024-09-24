package restic

import (
	"context"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindsnapshotsLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	commonRestic *commonsvc.CommonRestic
}

func NewFindsnapshotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindsnapshotsLogic {
	return &FindsnapshotsLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *FindsnapshotsLogic) Findsnapshots(req *types.FindSnapshotsReq) (resp *types.Response, err error) {
	var output string
	var args = []string{"find", req.Pattern, "--host", req.Host}
	if req.SnapshotId != "" {
		args = append(args, "-s", req.SnapshotId)
	}
	if req.Path != "" {
		args = append(args, "--path", req.Path)
	}
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
