package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hongyuxuan/lizardrestic/agent/internal/svc"
	"github.com/hongyuxuan/lizardrestic/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type BackupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	resticService *svc.ResticService
}

func NewBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BackupLogic {
	return &BackupLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		resticService: svc.NewResticService(ctx, svcCtx),
	}
}

func (l *BackupLogic) Backup(in *agent.BackupRequest, stream agent.LizardAgent_BackupServer) error {
	policyfile := fmt.Sprintf("%s/policy-%d", l.svcCtx.Config.ConfigurationDir, in.PolicyId)
	args := []string{"backup", "--files-from", policyfile, "--host", in.Host, "-v=1"}
	if len(in.Tags) > 0 {
		args = append(args, "--tag", strings.Join(in.Tags, ","))
	}
	if in.Exclude != "" {
		args = append(args, "--exclude", in.Exclude)
	}
	err := l.resticService.RunCommand(stream, in.RepoUrl, args...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
