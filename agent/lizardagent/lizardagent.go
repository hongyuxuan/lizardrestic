// Code generated by goctl. DO NOT EDIT!
// Source: agent.proto

package lizardagent

import (
	"context"

	"github.com/hongyuxuan/lizardrestic/agent/types/agent"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BackupRequest         = agent.BackupRequest
	CreatePolicyRequest   = agent.CreatePolicyRequest
	Response              = agent.Response
	RestoreRequest        = agent.RestoreRequest
	SetEnvironmentRequest = agent.SetEnvironmentRequest

	LizardAgent interface {
		SetEnvironment(ctx context.Context, in *SetEnvironmentRequest, opts ...grpc.CallOption) (*Response, error)
		CreatePolicy(ctx context.Context, in *CreatePolicyRequest, opts ...grpc.CallOption) (*Response, error)
		Backup(ctx context.Context, in *BackupRequest, opts ...grpc.CallOption) (agent.LizardAgent_BackupClient, error)
		Restore(ctx context.Context, in *RestoreRequest, opts ...grpc.CallOption) (agent.LizardAgent_RestoreClient, error)
	}

	defaultLizardAgent struct {
		cli zrpc.Client
	}
)

func NewLizardAgent(cli zrpc.Client) LizardAgent {
	return &defaultLizardAgent{
		cli: cli,
	}
}

func (m *defaultLizardAgent) SetEnvironment(ctx context.Context, in *SetEnvironmentRequest, opts ...grpc.CallOption) (*Response, error) {
	client := agent.NewLizardAgentClient(m.cli.Conn())
	return client.SetEnvironment(ctx, in, opts...)
}

func (m *defaultLizardAgent) CreatePolicy(ctx context.Context, in *CreatePolicyRequest, opts ...grpc.CallOption) (*Response, error) {
	client := agent.NewLizardAgentClient(m.cli.Conn())
	return client.CreatePolicy(ctx, in, opts...)
}

func (m *defaultLizardAgent) Backup(ctx context.Context, in *BackupRequest, opts ...grpc.CallOption) (agent.LizardAgent_BackupClient, error) {
	client := agent.NewLizardAgentClient(m.cli.Conn())
	return client.Backup(ctx, in, opts...)
}

func (m *defaultLizardAgent) Restore(ctx context.Context, in *RestoreRequest, opts ...grpc.CallOption) (agent.LizardAgent_RestoreClient, error) {
	client := agent.NewLizardAgentClient(m.cli.Conn())
	return client.Restore(ctx, in, opts...)
}