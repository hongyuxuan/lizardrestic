package restic

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletepolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletepolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletepolicyLogic {
	return &DeletepolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletepolicyLogic) Deletepolicy(req *types.BackupReq) (resp *types.Response, err error) {
	policy := &commontypes.BackupPolicy{}
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetBackupPolicy")).First(&policy, "id = ?", req.PolicyId).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	// delete cronjob
	svc.RemoveCron(l.svcCtx, *policy)
	// delete policy
	l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteBackupPolicy")).Delete(&commontypes.BackupPolicy{}, req.PolicyId)
	// delete policy history
	l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteBackupPolicy")).Delete(&commontypes.BackupHistory{}, "backup_policy_id = ?", req.PolicyId)

	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "删除策略成功",
	}
	return
}
