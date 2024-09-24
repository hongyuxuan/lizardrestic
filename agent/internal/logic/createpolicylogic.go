package logic

import (
	"context"
	"fmt"
	"os"

	"github.com/hongyuxuan/lizardrestic/agent/internal/svc"
	"github.com/hongyuxuan/lizardrestic/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePolicyLogic {
	return &CreatePolicyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePolicyLogic) CreatePolicy(in *agent.CreatePolicyRequest) (*agent.Response, error) {
	policyfile := fmt.Sprintf("%s/policy-%d", l.svcCtx.Config.ConfigurationDir, in.PolicyId)
	f, err := os.OpenFile(policyfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logx.Errorf("Failed to open file \"%s\": %v", policyfile, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	for _, line := range in.BackupDir {
		if _, err = f.WriteString(line + "\n"); err != nil {
			logx.Errorf("Failed to save policy file \"%s\": %v", policyfile, err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
