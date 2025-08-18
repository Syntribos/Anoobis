package main

import (
	"flag"
	"fmt"
	"os"
)

type param struct {
	name         string
	value        *string
	defaultValue string
	usage        string
	required     bool
}

const (
	tk = "token"
	rg = "reportguild"
	rc = "reportchannel"
	db = "dbpath"

	tkHelp = "Bot authorization token"
	rgHelp = "ID of server for user reports"
	rcHelp = "ID of channel for user reports"
	dbHelp = `Full path to database file.
	NOTE: database doesn't need to exist, but user must have permissions for the destination.`
)

var params = map[string]param{
	tk: param{tk, nil, "", tkHelp, true},
	rg: param{rg, nil, "", rgHelp, true},
	rc: param{rc, nil, "", rcHelp, true},
	db: param{db, nil, "", dbHelp, true},
}

func registerParams() {
	for key, value := range params {
		valuePtr := flag.String(value.name, value.defaultValue, value.usage)
		params[key] = param{value.name, valuePtr, value.defaultValue, value.usage, value.required}
	}
}

func validateParams() {
	valid := true
	for _, value := range params {
		if !value.required {
			continue
		}

		if value.value == nil || len(*value.value) == 0 {
			fmt.Println("Missing required parameter: " + value.name)
			valid = false
		}
	}

	if !valid {
		os.Exit(1)
	}
}
