package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetdataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetdataLogic {
	return &GetdataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetdataLogic) Getdata(req *types.DataByIdReq) (resp *types.Response, err error) {
	if req.Tablename == "backup_policy" {
		var policy commontypes.BackupPolicy
		if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetBackupPolicy")).
			First(&policy, "id = ?", req.Id).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetRepository")).
			First(&policy.Repository, "id = ?", policy.RepositoryId)
		resp = &types.Response{
			Code: http.StatusOK,
			Data: policy,
		}
	}
	return
}
