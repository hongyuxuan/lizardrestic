package logic

import (
	"context"
	"strings"

	"github.com/hongyuxuan/lizardrestic/agent/internal/svc"
	"github.com/hongyuxuan/lizardrestic/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	resticService *svc.ResticService
}

func NewRestoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreLogic {
	return &RestoreLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		resticService: svc.NewResticService(ctx, svcCtx),
	}
}

func (l *RestoreLogic) Restore(in *agent.RestoreRequest, stream agent.LizardAgent_RestoreServer) error {
	args := []string{"restore", in.SnapshotId, "--host", in.Host, "--target", in.Target, "-v=1"}
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
