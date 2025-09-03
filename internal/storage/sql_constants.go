package storage

import (
	_ "embed"
	m "github.com/Syntribos/Anoobis/internal/models"
)

const (
	currentMajor      int    = 0
	currentMinor      int    = 1
	getVersionCommand string = "SELECT major, minor FROM version"
	getReportCursor   string = "SELECT last_checked_message_id FROM historical_reports"
	saveReportCursor  string = "INSERT INTO historical_reports VALUES (0, ?)"
	saveReport        string = "INSERT INTO reports (message_id, report_link, user_id, reason) VALUES (?, ?, ?, ?)"
)

var currentVersion = m.DBVersion{Major: currentMajor, Minor: currentMinor}

//go:embed migrations/anoobis_info.sqlite.sql
var createCommand string
