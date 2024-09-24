package restic

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BackupLogic {
	return &BackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BackupLogic) Backup(req *types.BackupReq) (resp *types.Response, err error) {
	var policy commontypes.BackupPolicy
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetBackupPolicy")).
		Joins("Repository").
		First(&policy, "backup_policy.id = ?", req.PolicyId).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	for _, host := range policy.Hosts {
		go svc.DoBackup(l.svcCtx, policy, policy.Repository, host)
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "备份任务已提交",
	}
	return
}
