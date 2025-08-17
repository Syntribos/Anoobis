package storage

import (
	_ "embed"
	"fmt"
)

type DBInfo struct {
	dataSourceName string
	version        *DBVersion
}

type DBVersion struct {
	major int
	minor int
}

const (
	currentMajor      int    = 0
	currentMinor      int    = 1
	getVersionCommand string = "SELECT major, minor FROM version"
)

var currentVersion = DBVersion{currentMajor, currentMinor}

//go:embed migrations/anoobis_info.sqlite.sql
var createCommand string

func (dbVersion *DBVersion) GetVersionString() string {
	return fmt.Sprintf("%d.%d", dbVersion.major, dbVersion.minor)
}

func (dbInfo *DBInfo) GetVersionString() string {
	return dbInfo.version.GetVersionString()
}
