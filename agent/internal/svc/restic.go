package svc

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"

	"github.com/hongyuxuan/lizardrestic/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"
)

type ResticService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func NewResticService(ctx context.Context, svcCtx *ServiceContext) *ResticService {
	return &ResticService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (r *ResticService) RunCommand(stream agent.LizardAgent_BackupServer, repoUrl string, args ...string) (err error) {
	envfile := fmt.Sprintf("%s/environment-%s.yaml", r.svcCtx.Config.ConfigurationDir, url.QueryEscape(repoUrl))
	var envYaml []byte
	if envYaml, err = os.ReadFile(envfile); err != nil {
		r.Logger.Error(err)
		return
	}
	var env commontypes.ResticEnvironment
	if err = yaml.Unmarshal(envYaml, &env); err != nil {
		r.Logger.Error(err)
		return
	}
	args = append([]string{"--cache-dir", r.svcCtx.Config.CacheDir}, args...)
	cmd := exec.Command("restic", args...)
	cmd.Env = []string{
		"RESTIC_REPOSITORY=" + env.RESTIC_REPOSITORY,
		"RESTIC_PASSWORD=" + env.RESTIC_PASSWORD,
		"AWS_ACCESS_KEY_ID=" + env.AWS_ACCESS_KEY_ID,
		"AWS_SECRET_ACCESS_KEY=" + env.AWS_SECRET_ACCESS_KEY,
	}
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		r.Logger.Error(err)
		return
	}
	cmd.Stderr = cmd.Stdout
	if err = cmd.Start(); err != nil {
		r.Logger.Error(err)
		return
	}
	for {
		b := make([]byte, 1024)
		var n int
		n, err = stdout.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			r.Logger.Error(err)
			return
		}
		output := string(b[:n])
		r.Logger.Info(output)
		stream.Send(&agent.Response{
			Code:    uint32(codes.OK),
			Message: output,
		})
	}
	if err = cmd.Wait(); err != nil {
		r.Logger.Error(err)
		return
	}
	return nil
}
