package storage

import (
	"database/sql"
	"errors"
	m "github.com/Syntribos/Anoobis/internal/models"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

func Init(fullPath string) (*m.DBInfo, error) {
	var err error
	var version *m.DBVersion
	version, err = GetVersion(fullPath)

	if errors.Is(err, os.ErrNotExist) {
		log.Println("DB doesn't exist at " + fullPath)
		version, err = CreateDatabase(fullPath)
	}

	if err != nil {
		return nil, err
	}

	if version.Major != currentVersion.Major ||
		version.Minor != currentVersion.Minor {
		// Upgrade schema
	}

	d := &m.DBInfo{
		DataSourceName: fullPath,
		Version:        version,
	}

	return d, nil
}

func GetVersion(fullPath string) (*m.DBVersion, error) {
	var db *sql.DB
	var major, minor int

	if _, err := os.Stat(fullPath); err == nil {
		db, err = sql.Open("sqlite", fullPath)

		if err != nil {
			return nil, err
		}

		row := db.QueryRow(
			getVersionCommand,
		)

		if err = row.Scan(&major, &minor); err != nil {
			return nil, err
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// db doesn't exist, raise error to caller
		return nil, err

	} else {
		// secret third option, check err

	}

	v := &m.DBVersion{Major: major, Minor: minor}

	return v, nil
}

func CreateDatabase(fullPath string) (*m.DBVersion, error) {
	if len(createCommand) == 0 {
		return nil, errors.New("CreateDatabase command missing")
	}

	log.Println("Creating database at " + fullPath)
	db, err := sql.Open("sqlite", fullPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createCommand); err != nil {
		log.Fatal(err)
	}

	log.Println("Database created at " + fullPath)

	return GetVersion(fullPath)
}

func GetCurrentReportCursor(dbInfo *m.DBInfo, channelId string) (string, error) {
	db, err := sql.Open("sqlite", dbInfo.DataSourceName)
	result := ""

	defer func() {
		err = db.Close()
	}()

	if err != nil {
		return result, err
	}

	row := db.QueryRow(
		getReportCursor,
	)

	if err = row.Scan(&result); err != nil {
		return "", err
	}

	return result, err
}

func SaveReportCursor(dbInfo *m.DBInfo, messageId string) error {
	db, err := sql.Open("sqlite", dbInfo.DataSourceName)

	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	_, err = db.Exec(saveReportCursor, messageId)

	return err
}

func SaveReport(dbInfo *m.DBInfo, report m.ReportPackage) error {
	db, err := sql.Open("sqlite", dbInfo.DataSourceName)
	if err != nil {
		return err
	}

	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	statement, err := transaction.Prepare(saveReport)
	if err != nil {
		txErr := transaction.Rollback()
		if txErr != nil {
			log.Fatal("Unable to rollback transaction:", txErr)
		}
		return err
	}

	for i := range report.Reports {
		currReport := &report.Reports[i]
		reportLink := report.GetMessageLink(i)
		_, err = statement.Exec(currReport.MessageId, reportLink, currReport.UserId, report.Reason)
		if err != nil {
			txErr := transaction.Rollback()
			if txErr != nil {
				log.Fatal("Unable to rollback transaction:", txErr)
			}
			return err
		}
	}

	return transaction.Commit()
}
