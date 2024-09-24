package restic

import (
	"context"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreaterepoLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	commonRestic *commonsvc.CommonRestic
}

func NewCreaterepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreaterepoLogic {
	return &CreaterepoLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *CreaterepoLogic) Createrepo(req *types.CreateRepoReq) (resp *types.Response, err error) {
	repo := commontypes.Repository{
		RepoName:    req.RepoName,
		RepoUrl:     req.RepoUrl,
		S3AccessKey: req.S3AccessKey,
		S3SecretKey: req.S3SecretKey,
		Password:    req.Password,
	}
	// insert database
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.SaveRepository")).
		Save(&repo).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	// set server environment
	if err = l.commonRestic.SetEnvironment(req.RepoUrl, req.Password, req.S3AccessKey, req.S3SecretKey); err != nil {
		return
	}
	// init repository
	var output string
	if output, err = l.commonRestic.RunCommand(req.RepoUrl, "init"); err != nil {
		return
	}
	l.Logger.Infof("Init repository %s success: %s", req.RepoUrl, output)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: output,
	}
	return
}
