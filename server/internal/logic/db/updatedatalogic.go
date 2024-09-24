package db

import (
	"context"

	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatedataLogic {
	return &UpdatedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatedataLogic) Updatedata(req *types.UpdateDataReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
