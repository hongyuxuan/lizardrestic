package svc

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/hongyuxuan/lizardrestic/common/errorx"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"
)

type CommonRestic struct {
	logx.Logger
	ctx              context.Context
	configurationDir string
	cacheDir         string
}

func NewCommonRestic(ctx context.Context, configurationDir, cacheDir string) *CommonRestic {
	return &CommonRestic{
		Logger:           logx.WithContext(ctx),
		ctx:              ctx,
		configurationDir: configurationDir,
		cacheDir:         cacheDir,
	}
}

func (r *CommonRestic) SetEnvironment(repoUrl, password, accessKey, secretKey string) (err error) {
	envfile := fmt.Sprintf("%s/environment-%s.yaml", r.configurationDir, url.QueryEscape(repoUrl))
	f, err := os.OpenFile(envfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logx.Errorf("Failed to open file \"%s\": %v", envfile, err)
		return
	}
	env := commontypes.ResticEnvironment{
		RESTIC_REPOSITORY:     repoUrl,
		RESTIC_PASSWORD:       password,
		AWS_ACCESS_KEY_ID:     accessKey,
		AWS_SECRET_ACCESS_KEY: secretKey,
	}
	envYaml, _ := yaml.Marshal(env)
	if _, err = f.Write(envYaml); err != nil {
		logx.Errorf("Failed to save environment into file \"%s\": %v", envfile, err)
		return
	}
	logx.Infof("Save environment into file \"%s\" success", envfile)
	return
}

func (r *CommonRestic) RunCommand(repoUrl string, args ...string) (output string, err error) {
	envfile := fmt.Sprintf("%s/environment-%s.yaml", r.configurationDir, url.QueryEscape(repoUrl))
	var envYaml []byte
	if envYaml, err = os.ReadFile(envfile); err != nil {
		logx.Error(err)
		return
	}
	var env commontypes.ResticEnvironment
	if err = yaml.Unmarshal(envYaml, &env); err != nil {
		logx.Error(err)
		return
	}
	args = append([]string{"--cache-dir", r.cacheDir}, args...)
	cmd := exec.Command("restic", args...)
	cmd.Env = []string{
		"RESTIC_REPOSITORY=" + env.RESTIC_REPOSITORY,
		"RESTIC_PASSWORD=" + env.RESTIC_PASSWORD,
		"AWS_ACCESS_KEY_ID=" + env.AWS_ACCESS_KEY_ID,
		"AWS_SECRET_ACCESS_KEY=" + env.AWS_SECRET_ACCESS_KEY,
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		err = errorx.NewDefaultError(err.Error() + "; " + out.String())
		logx.Error(err)
		return
	}
	return out.String(), nil
}
