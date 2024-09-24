package restic

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreLogic {
	return &RestoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestoreLogic) Restore(req *types.RestoreReq) (resp *types.Response, err error) {
	var policy commontypes.BackupPolicy
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetBackupPolicy")).
		Joins("Repository").
		First(&policy, "backup_policy.id = ?", req.PolicyId).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	go svc.DoRestore(l.svcCtx, policy, policy.Repository, req.SnapshotId, req.Host, req.Exclude, req.Target)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "恢复任务已提交",
	}
	return
}
