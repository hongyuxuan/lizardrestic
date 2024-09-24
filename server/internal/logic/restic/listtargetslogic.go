package restic

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardrestic/common/utils"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListtargetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListtargetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListtargetsLogic {
	return &ListtargetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListtargetsLogic) Listtargets() (resp *types.Response, err error) {
	var targets []map[string]string
	for k := range l.svcCtx.AgentList {
		system, ip, err := utils.GetTarget(l.svcCtx.Config.ServicePrefix, k)
		if err != nil {
			continue
		}
		targets = append(targets, map[string]string{"system": system, "ip": ip})
	}
	if targets == nil {
		targets = []map[string]string{}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: targets,
	}
	return
}
