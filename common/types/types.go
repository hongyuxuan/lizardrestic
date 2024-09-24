package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type TraceIDKey struct{}

type ResticEnvironment struct {
	RESTIC_REPOSITORY     string
	RESTIC_PASSWORD       string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "[]"
	}
	return str, err
}

func (s *StringList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "{}"
	}
	return str, err
}

func (s *StringMap) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

type GetDataReq struct {
	Tablename string `path:"tablename"`
	Page      int    `form:"page,default=1"`
	Size      int    `form:"size,default=20"`
	Search    string `form:"search,optional"`
	Filter    string `form:"filter,optional"`
	Range     string `form:"range,optional"`
	Preload   bool   `form:"preload,optional"`
	Sort      string `form:"sort,optional"`
}

type Repository struct {
	Id            int    `json:"id" gorm:"primaryKey,autoIncrement"`
	RepoName      string `json:"repo_name" gorm:"size:50;unique"`
	RepoUrl       string `json:"repo_url" gorm:"size:300"`
	S3AccessKey   string `json:"s3_access_key" gorm:"size:100"`
	S3SecretKey   string `json:"s3_secret_key" gorm:"size:100"`
	Password      string `json:"password" gorm:"size:100"`
	FilesCount    int    `json:"total_file_count"`
	TotalSize     string `json:"total_size" gorm:"size:50"`
	SnapshotCount int    `json:"snapshots_count"`
}

type BackupPolicy struct {
	Id           int        `json:"id" gorm:"primaryKey,autoIncrement"`
	Name         string     `json:"name" gorm:"size:100;unique"`
	RepositoryId int        `json:"repository_id"`
	Hosts        StringList `json:"hosts" gorm:"type:json"`
	Cron         string     `json:"cron" gorm:"size:100"`
	BackupDir    string     `json:"backup_dir"`
	Tags         StringList `json:"tags" gorm:"type:json"`
	Enable       bool       `json:"enable"`
	Exclude      string     `json:"exclude" gorm:"size:300"`
	Retention    string     `json:"retention" gorm:"size:50"`
	Repository   Repository `json:"repository,omitempty" gorm:"foreignKey:RepositoryId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type BackupHistory struct {
	Id             int          `json:"id" gorm:"primaryKey,autoIncrement"`
	BackupPolicyId int          `json:"repository_id"`
	BackupPolicy   BackupPolicy `json:"repository,omitempty" gorm:"foreignKey:BackupPolicyId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TaskType       string       `json:"task_type" gorm:"size:50"`
	Host           string       `json:"host" gorm:"size:50"                                 `
	Success        sql.NullBool `json:"success"`
	Message        string       `json:"message"`
	Status         string       `json:"status" gorm:"size:50"`
	StartAt        sql.NullTime `json:"start_at"`
	FinishAt       sql.NullTime `json:"finish_at"`
}

type ResticSnapshot struct {
	Id       string    `json:"short_id"`
	Time     time.Time `json:"time"`
	Paths    []string  `json:"paths"`
	Tags     []string  `json:"tags"`
	Hostname string    `json:"hostname"`
	Summary  struct {
		FilesNew        int64 `json:"files_new"`
		FilesChanged    int64 `json:"files_changed"`
		FilesUnmodified int64 `json:"files_unmodified"`
		SizeInt64       int64 `json:"total_bytes_processed"`
	} `json:"summary"`
	Size string `json:"size"`
}

type RepositoryStats struct {
	TotalSize      int64 `json:"total_size"`
	TotalFileCount int64 `json:"total_file_count"`
	SnapshotsCount int64 `json:"snapshots_count"`
}

type ListResult struct {
	Total   int         `json:"total"`
	Results interface{} `json:"results"`
}
