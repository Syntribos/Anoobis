package main

import (
	"flag"
	"fmt"
	"github.com/Syntribos/Anoobis/internal/anoobis_client"
)

var (
	BotToken = flag.String("token", "", "Bot authorization token")
	GuildId  = flag.String("guild", "", "ID of the testing guild")
	DbPath   = flag.String("dbpath", "", "Full path to the database file")
)

func init() { flag.Parse() }

func main() {
	valid := validateParam("token", BotToken)
	valid = valid && validateParam("dbpath", DbPath)

	if !valid {
		fmt.Println("Please rerun with valid parameters.")
		return
	}

	anoobis_client.Run(BotToken, GuildId, DbPath)
}

func validateParam(name string, value *string) bool {
	if value == nil || len(*value) == 0 {
		fmt.Println("Missing required value for '" + name + "'")
		return false
	}
	return true
}
