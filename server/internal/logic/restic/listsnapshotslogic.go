package restic

import (
	"context"
	"encoding/json"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/common/utils"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListsnapshotsLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	commonRestic *commonsvc.CommonRestic
}

func NewListsnapshotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListsnapshotsLogic {
	return &ListsnapshotsLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *ListsnapshotsLogic) Listsnapshots(req *types.ListSnapshotReq) (resp *types.Response, err error) {
	var policy commontypes.BackupPolicy
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListSnapshots")).
		Model(&commontypes.BackupPolicy{}).
		First(&policy, "id = ?", req.PolicyId).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	var repo commontypes.Repository
	l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetRepository")).
		Model(&commontypes.Repository{}).
		Find(&repo, "id = ?", policy.RepositoryId)
	var output string
	var args = []string{"snapshots", "--host", req.Host, "--json"}
	if req.Tag != "" {
		args = append(args, "--tag", req.Tag)
	}
	if req.Latest != "" {
		args = append(args, "--latest", req.Latest)
	}
	if output, err = l.commonRestic.RunCommand(repo.RepoUrl, args...); err != nil {
		return
	}
	var snapshots []*commontypes.ResticSnapshot
	if err = json.Unmarshal([]byte(output), &snapshots); err != nil {
		l.Logger.Error(err)
		return
	}
	for _, snapshot := range snapshots {
		snapshot.Size = utils.ByteCountIEC(snapshot.Summary.SizeInt64)
	}
	l.Logger.Infof("List %s snapshots success: total %d snapshots", repo.RepoUrl, len(snapshots))
	if len(snapshots) == 0 {
		snapshots = []*commontypes.ResticSnapshot{}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: snapshots,
	}
	return
}
