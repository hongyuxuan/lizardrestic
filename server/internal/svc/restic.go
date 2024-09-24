package svc

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/hongyuxuan/lizardrestic/agent/lizardagent"
	"github.com/hongyuxuan/lizardrestic/common/constant"
	commonsvc "github.com/hongyuxuan/lizardrestic/common/svc"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
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

func (r *ResticService) LoadCronJob() (err error) {
	var policies []commontypes.BackupPolicy
	if err = r.svcCtx.Sqlite.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.ListBackupPolicy")).
		Model(&commontypes.BackupPolicy{}).
		Joins("Repository").
		Where("enable = true").
		Find(&policies).Error; err != nil {
		r.Logger.Errorf("Failed to load backup policies from database: %v", err)
		return
	}
	for _, policy := range policies {
		AddCron(r.svcCtx, policy, policy.Repository)
	}
	r.svcCtx.Cron.Start()
	return
}

func AddCron(svcCtx *ServiceContext, policy commontypes.BackupPolicy, repo commontypes.Repository) (err error) {
	commonRestic := commonsvc.NewCommonRestic(context.Background(), svcCtx.Config.ConfigurationDir, svcCtx.Config.CacheDir)
	for _, host := range policy.Hosts {
		var cronId cron.EntryID
		if cronId, err = svcCtx.Cron.AddFunc(policy.Cron, func() {
			DoBackup(svcCtx, policy, repo, host)
			// delete snapshots due to retention
			if policy.Retention != "" {
				args := []string{"forget", "--host", host, "--keep-within", policy.Retention, "--prune"}
				if len(policy.Tags) > 0 {
					args = append(args, "--tag", strings.Join(policy.Tags, ","))
				}
				var output string
				if output, err = commonRestic.RunCommand(repo.RepoUrl, args...); err != nil {
					logx.Errorf("Failed to delete snapshots with host=%s tag=%s retention=%s", host, strings.Join(policy.Tags, ","), policy.Retention)
				} else {
					logx.Infof("Successfully delete snapshots with host=%s tag=%s retention=%s: %v", host, strings.Join(policy.Tags, ","), policy.Retention, output)
				}
			}
		}); err != nil {
			logx.Errorf("Failed to add cronjob: %v", err)
			continue
		}
		svcCtx.CronIdMap[strconv.Itoa(policy.Id)+"-"+host] = cronId
		logx.Infof("Add cronjob cron=%s host=%s backup_dir=%s success", policy.Cron, host, strings.Join(strings.Split(policy.BackupDir, "\n"), ","))
	}
	return
}

func RemoveCron(svcCtx *ServiceContext, policy commontypes.BackupPolicy) {
	for _, host := range policy.Hosts {
		svcCtx.Cron.Remove(svcCtx.CronIdMap[strconv.Itoa(policy.Id)+"-"+host])
		logx.Infof("Remove cronjob %d-%s success", policy.Id, host)
	}
}

func DoBackup(svcCtx *ServiceContext, policy commontypes.BackupPolicy, repo commontypes.Repository, host string) {
	// add to history
	history := commontypes.BackupHistory{
		BackupPolicyId: policy.Id,
		Host:           host,
		Message:        "",
		TaskType:       constant.TASK_TYPE_BACKUP,
		Status:         constant.CRON_STATUS_RUNNING,
		StartAt:        sql.NullTime{Time: time.Now(), Valid: true},
	}
	if err := svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.SaveBackupHistory")).
		Save(&history).Error; err != nil {
		logx.Error(err)
	}
	// send backup to agent
	ag, err := svcCtx.GetTargetAgent(host)
	if err != nil {
		svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.UpdateBackupHistory")).
			Exec("UPDATE backup_history SET success = ?, status = ?, finish_at = ?, message = ? WHERE id = ?", false, constant.CRON_STATUS_FINISHED, time.Now().String(), err.Error(), history.Id)
		return
	}
	stream, _ := ag.Backup(context.Background(), &lizardagent.BackupRequest{
		PolicyId: uint32(policy.Id),
		RepoUrl:  repo.RepoUrl,
		Host:     host,
		Tags:     policy.Tags,
	})
	for {
		resp, e := stream.Recv()
		if e == io.EOF {
			svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.SaveBackupHistory")).
				Model(&commontypes.BackupHistory{}).
				Where("id = ?", history.Id).
				Updates(&commontypes.BackupHistory{
					Success:  sql.NullBool{Bool: true, Valid: true},
					Status:   constant.CRON_STATUS_FINISHED,
					FinishAt: sql.NullTime{Time: time.Now(), Valid: true},
				})
			break
		}
		if e != nil {
			svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.UpdateBackupHistory")).
				Exec("UPDATE backup_history SET success = ?, status = ?, finish_at = ?, message = message || ? WHERE id = ?", false, constant.CRON_STATUS_FINISHED, time.Now().String(), e.Error(), history.Id)
			break
		}
		svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.UpdateBackupHistory")).
			Exec("UPDATE backup_history SET status = ?, message = message || ? || '\n' WHERE id = ?", constant.CRON_STATUS_RUNNING, resp.Message, history.Id)
		logx.Debugf("Cronjob running, host=%s recv message=%s", host, resp.Message)
	}
	logx.Infof("Execute crojob cron=%s host=%s backup_dir=%s success", policy.Cron, host, strings.Join(strings.Split(policy.BackupDir, "\n"), ","))
	return
}

func DoRestore(svcCtx *ServiceContext, policy commontypes.BackupPolicy, repo commontypes.Repository, snapshotId, host, exclude, targets string) {
	// add to history
	history := commontypes.BackupHistory{
		BackupPolicyId: policy.Id,
		Host:           host,
		Message:        "",
		TaskType:       constant.TASK_TYPE_RESTORE,
		Status:         constant.CRON_STATUS_RUNNING,
		StartAt:        sql.NullTime{Time: time.Now(), Valid: true},
	}
	if err := svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.SaveBackupHistory")).
		Save(&history).Error; err != nil {
		logx.Error(err)
	}
	// send backup to agent
	ag, err := svcCtx.GetTargetAgent(host)
	if err != nil {
		svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.UpdateBackupHistory")).
			Exec("UPDATE backup_history SET success = ?, status = ?, finish_at = ?, message = ? WHERE id = ?", false, constant.CRON_STATUS_FINISHED, time.Now().String(), err.Error(), history.Id)
		return
	}
	stream, _ := ag.Restore(context.Background(), &lizardagent.RestoreRequest{
		RepoUrl:    repo.RepoUrl,
		SnapshotId: snapshotId,
		Host:       host,
		Tags:       policy.Tags,
		Exclude:    exclude,
		Target:     targets,
	})
	for {
		resp, e := stream.Recv()
		if e == io.EOF {
			svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.SaveBackupHistory")).
				Model(&commontypes.BackupHistory{}).
				Where("id = ?", history.Id).
				Updates(&commontypes.BackupHistory{
					Success:  sql.NullBool{Bool: true, Valid: true},
					Status:   constant.CRON_STATUS_FINISHED,
					FinishAt: sql.NullTime{Time: time.Now(), Valid: true},
				})
			break
		}
		if e != nil {
			svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.UpdateBackupHistory")).
				Exec("UPDATE backup_history SET success = ?, status = ?, finish_at = ?, message = message || ? WHERE id = ?", false, constant.CRON_STATUS_FINISHED, time.Now().String(), e.Error(), history.Id)
			break
		}
		svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.UpdateBackupHistory")).
			Exec("UPDATE backup_history SET status = ?, message = message || ? || '\n' WHERE id = ?", constant.CRON_STATUS_RUNNING, resp.Message, history.Id)
		logx.Debugf("Cronjob running, host=%s recv message=%s", host, resp.Message)
	}
	logx.Infof("Execute crojob cron=%s host=%s backup_dir=%s success", policy.Cron, host, strings.Join(strings.Split(policy.BackupDir, "\n"), ","))
	return
}
