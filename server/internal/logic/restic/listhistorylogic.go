package restic

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListhistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListhistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListhistoryLogic {
	return &ListhistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListhistoryLogic) Listhistory(req *types.ListHistoryReq) (resp *types.Response, err error) {
	var res []map[string]interface{}
	var count map[string]interface{}
	tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListHistory")).
		Model(&commontypes.BackupHistory{}).
		Select("backup_policy.name policy_name,backup_history.*").
		Joins("left join backup_policy on backup_history.backup_policy_id=backup_policy.id")
	txCount := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListHistory")).
		Model(&commontypes.BackupHistory{}).
		Select("backup_policy.name policy_name,count(*) count").
		Joins("left join backup_policy on backup_history.backup_policy_id=backup_policy.id")
	if req.Search != "" {
		for _, search := range strings.Split(req.Search, ",") {
			searchStmt := strings.Split(search, "==")
			tx.Where(fmt.Sprintf("%s LIKE ?", searchStmt[0]), "%"+searchStmt[1]+"%")
			txCount.Where(fmt.Sprintf("%s LIKE ?", searchStmt[0]), "%"+searchStmt[1]+"%")
		}
	}
	if req.PolicyName != "" {
		tx.Where("policy_name LIKE ?", "%"+req.PolicyName+"%")
		txCount.Where("policy_name LIKE ?", "%"+req.PolicyName+"%")
	}
	if req.Filter != "" {
		for _, filter := range strings.Split(req.Filter, ",") {
			filterStmt := strings.Split(filter, "==")
			tx.Where(fmt.Sprintf("%s = ?", filterStmt[0]), filterStmt[1])
			txCount.Where(fmt.Sprintf("%s = ?", filterStmt[0]), filterStmt[1])
		}
	}
	if req.Range != "" { // &range=init_at==2024-02-06 16:58:30,2024-02-06 17:58:30
		rangeS := strings.Split(req.Range, "==")
		rangeKey := rangeS[0]
		rangeR := strings.Split(rangeS[1], ",")
		tx.Where(fmt.Sprintf("%s BETWEEN ? AND ?", rangeKey), rangeR[0], rangeR[1])
		txCount.Where(fmt.Sprintf("%s BETWEEN ? AND ?", rangeKey), rangeR[0], rangeR[1])
	}
	if err = txCount.First(&count).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	tx.Limit(req.Size).Offset((req.Page - 1) * req.Size)
	if req.Sort != "" {
		tx.Order(req.Sort)
	}
	if err = tx.Find(&res).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count["count"].(int64)),
			Results: res,
		},
	}
	return
}
