package main

import (
	"github.com/alecthomas/kingpin/v2"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/common/utils"
	"go.opentelemetry.io/otel"
)

var (
	level     = kingpin.Flag("log.level", "Log level.").Default("info").String()
	dbfile    = kingpin.Flag("db", "SQLite database file.").Required().String()
	serverUrl = kingpin.Flag("server-url", "LizardRestic server url").Short('s').Default("http://localhost:7138").String()
)

type RepoStatsResponse struct {
	Code int                         `json:"code"`
	Data commontypes.RepositoryStats `json:"data"`
}

func main() {
	// Parse flags
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	db := utils.NewSQLite(*dbfile, *level)
	http := utils.NewHttpClient(otel.Tracer("imroc/req"))
	http.SetBaseURL(*serverUrl)

	var repos []commontypes.Repository
	if err := db.Find(&repos).Error; err != nil {
		utils.Log.Fatal(err)
	}
	for _, repo := range repos {
		var res *RepoStatsResponse
		if err := http.Get("/restic/repository/stats").SetQueryParam("repo_url", repo.RepoUrl).SetSuccessResult(&res).Do().Err; err != nil {
			utils.Log.Error(err)
			continue
		}
		db.Model(&commontypes.Repository{}).Where("id = ?", repo.Id).Updates(&commontypes.Repository{
			FilesCount:    int(res.Data.TotalFileCount),
			TotalSize:     utils.ByteCountIEC(res.Data.TotalSize),
			SnapshotCount: int(res.Data.SnapshotsCount),
		})
	}
}
