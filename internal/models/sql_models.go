package models

import "fmt"

type DBInfo struct {
	DataSourceName string
	Version        *DBVersion
}

type DBVersion struct {
	Major int
	Minor int
}

func (dbVersion *DBVersion) GetVersionString() string {
	return fmt.Sprintf("%d.%d", dbVersion.Major, dbVersion.Minor)
}

func (dbInfo *DBInfo) GetVersionString() string {
	return dbInfo.Version.GetVersionString()
}
