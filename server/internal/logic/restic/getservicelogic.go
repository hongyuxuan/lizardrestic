package restic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hongyuxuan/lizardrestic/common/utils"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetserviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetserviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetserviceLogic {
	return &GetserviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetserviceLogic) Getservice(req *types.GetServiceReq) (resp *types.Response, err error) {
	var res *clientv3.GetResponse
	if res, err = l.svcCtx.EtcdClient.Get(l.ctx, req.ServiceName, clientv3.WithPrefix()); err != nil {
		logx.Error(err)
		return
	}
	var keymaps []map[string]interface{}
	for _, kv := range res.Kvs {
		key := utils.GetLizardAgentKey(kv.Key)
		meta, _ := utils.GetServiceMata(l.svcCtx.Config.ServicePrefix, key)
		keymaps = append(keymaps, map[string]interface{}{
			"ServiceID":   fmt.Sprintf("%s-%s", key, string(kv.Value)),
			"ServiceName": key,
			"ServiceMeta": meta,
		})
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: keymaps,
	}
	return
}
