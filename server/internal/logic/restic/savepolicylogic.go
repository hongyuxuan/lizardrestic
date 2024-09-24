package restic

import (
	"context"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardrestic/agent/lizardagent"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type SavepolicyLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	resticService *svc.ResticService
}

func NewSavepolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavepolicyLogic {
	return &SavepolicyLogic{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		svcCtx:        svcCtx,
		resticService: svc.NewResticService(ctx, svcCtx),
	}
}

func (l *SavepolicyLogic) Savepolicy(req *types.SavePolicyReq) (resp *types.Response, err error) {
	var repo *commontypes.Repository
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetRepository")).
		Where("id = ?", req.RepositoryId).
		First(&repo).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	for _, host := range req.Hosts {
		var ag lizardagent.LizardAgent
		if ag, err = l.svcCtx.GetTargetAgent(host); err != nil {
			return
		}
		// set agent environment
		if _, err = ag.SetEnvironment(l.ctx, &lizardagent.SetEnvironmentRequest{
			RepoUrl:     repo.RepoUrl,
			Password:    repo.Password,
			S3AccessKey: repo.S3AccessKey,
			S3SecretKey: repo.S3SecretKey,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	// save policy to database
	policy := &commontypes.BackupPolicy{}
	if err = copier.Copy(&policy, &req); err != nil {
		l.Logger.Error(err)
		return
	}
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.SaveBackupPolicy")).Save(&policy).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	// save policy file on agent
	for _, host := range req.Hosts {
		var ag lizardagent.LizardAgent
		if ag, err = l.svcCtx.GetTargetAgent(host); err != nil {
			return
		}
		if _, err = ag.CreatePolicy(l.ctx, &lizardagent.CreatePolicyRequest{
			PolicyId:  uint32(policy.Id),
			BackupDir: strings.Split(req.BackupDir, "\n"),
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	// add cronjob
	if req.Id == 0 && req.Enable {
		if err = svc.AddCron(l.svcCtx, *policy, *repo); err != nil {
			l.Logger.Error(err)
		}
	}
	if req.Id != 0 {
		svc.RemoveCron(l.svcCtx, *policy)
		if req.Enable {
			if err = svc.AddCron(l.svcCtx, *policy, *repo); err != nil {
				l.Logger.Error(err)
			}
		}
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "保存策略成功",
	}
	return
}
