package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedataLogic {
	return &CreatedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatedataLogic) Createdata(req *types.CreateDataReq) (resp *types.Response, err error) {
	data := make(map[string]interface{})
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.CreateData")).
		Table(req.Tablename).
		Create(&req.Body).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "新增成功",
		Data:    data,
	}
	return
}
