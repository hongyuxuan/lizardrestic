package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/common/utils"
)

var (
	level  = kingpin.Flag("log.level", "Log level.").Default("").String()
	dbfile = kingpin.Flag("db", "SQLite database file.").Short('d').Default("").String()
)

func main() {
	// Parse flags
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	db := utils.NewSQLite(*dbfile, *level)

	// create table `repository``
	if err := db.AutoMigrate(&types.Repository{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `backup_policy``
	if err := db.AutoMigrate(&types.BackupPolicy{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `backup_history``
	if err := db.AutoMigrate(&types.BackupHistory{}); err != nil {
		utils.Log.Warn(err)
	}
}
