package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/common/utils"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListdataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListdataLogic {
	return &ListdataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListdataLogic) Listdata(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	if req.Tablename == "backup_policy" {
		var data []commontypes.BackupPolicy
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListBackupPolicy")).Model(commontypes.BackupPolicy{}), data, req)
	}
	// []map[string]interface{} cannot use l.list, will be failed with 'sql: Scan error on column index 0, name "id": destination not a pointer'
	var data []map[string]interface{}
	tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListData")).Table(req.Tablename)
	var count int64
	utils.SetTx(tx, &count, req)
	if err = tx.Find(&data).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	if data == nil {
		data = []map[string]interface{}{}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count),
			Results: data,
		},
	}
	return
}

func (l *ListdataLogic) list(tx *gorm.DB, models any, req *commontypes.GetDataReq) (resp *types.Response, err error) {
	var count int64
	utils.SetTx(tx, &count, req)
	if err = tx.Find(&models).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count),
			Results: models,
		},
	}
	return
}
