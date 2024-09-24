package restic

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListcronsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListcronsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListcronsLogic {
	return &ListcronsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListcronsLogic) Listcrons() (resp *types.Response, err error) {
	resp = &types.Response{
		Code: http.StatusOK,
		Data: l.svcCtx.CronIdMap,
	}
	return
}
