package restic

import (
	"context"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesnapshotsLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	commonRestic *commonsvc.CommonRestic
}

func NewDeletesnapshotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesnapshotsLogic {
	return &DeletesnapshotsLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *DeletesnapshotsLogic) Deletesnapshots(req *types.SnapshotsIdsReq) (resp *types.Response, err error) {
	var output string
	ids := append([]string{"forget"}, req.Ids...)
	if output, err = l.commonRestic.RunCommand(req.RepoUrl, ids...); err != nil {
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: output,
	}
	return
}
