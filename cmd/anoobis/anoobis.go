package main

import (
	"flag"
	"fmt"
	"github.com/Syntribos/Anoobis/internal/anoobis_client"
	"log"
	"os"
)

func init() {
	registerParams()
	flag.Parse()
	validateParams()
}

func main() {
	var err error

	token := params[tk].value
	guildId := params[rg].value
	dbPath := params[db].value

	if err = anoobis_client.Run(token, guildId, dbPath); err == nil {
		os.Exit(0)
	}

	log.Fatal(err)
}

func validateParam(name string, value *string) bool {
	if value == nil || len(*value) == 0 {
		fmt.Println("Missing required value for '" + name + "'")
		return false
	}
	return true
}
