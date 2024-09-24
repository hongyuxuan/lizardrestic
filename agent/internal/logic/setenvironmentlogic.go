package logic

import (
	"context"

	"github.com/hongyuxuan/lizardrestic/agent/internal/svc"
	"github.com/hongyuxuan/lizardrestic/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetEnvironmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	commonRestic *commonsvc.CommonRestic
}

func NewSetEnvironmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetEnvironmentLogic {
	return &SetEnvironmentLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		commonRestic: commonsvc.NewCommonRestic(ctx, svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir),
	}
}

func (l *SetEnvironmentLogic) SetEnvironment(in *agent.SetEnvironmentRequest) (*agent.Response, error) {
	err := l.commonRestic.SetEnvironment(in.RepoUrl, in.Password, in.S3AccessKey, in.S3SecretKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
