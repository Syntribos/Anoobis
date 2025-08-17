package storage

import (
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
	"os"
)

func Init(fullPath string) (*DBInfo, error) {
	var err error
	var version *DBVersion
	version, err = GetVersion(fullPath)

	if errors.Is(err, os.ErrNotExist) {
		version, err = CreateDatabase(fullPath)
	}

	if err != nil {
		return nil, err
	}

	if version.major != currentVersion.major ||
		version.minor != currentVersion.minor {
		// Upgrade schema
	}

	d := &DBInfo{
		dataSourceName: fullPath,
		version:        version,
	}

	return d, nil
}

func GetVersion(fullPath string) (*DBVersion, error) {
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

		err = row.Scan(&major, &minor)

		if err != nil {
			return nil, err
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// db doesn't exist, raise error to caller
		return nil, err

	} else {
		// secret third option, check err

	}

	v := &DBVersion{major, minor}

	return v, nil
}

func CreateDatabase(fullPath string) (*DBVersion, error) {
	if len(createCommand) == 0 {
		return nil, errors.New("CreateDatabase command missing")
	}

	db, err := sql.Open("sqlite", fullPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createCommand); err != nil {
		return nil, err
	}

	return GetVersion(fullPath)
}
