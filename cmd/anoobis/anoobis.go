package main

import (
	"flag"
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

	if err = anoobis_client.Run(params[tk].value, params[rg].value, params[rc].value, params[db].value); err == nil {
		os.Exit(0)
	}

	log.Fatal(err)
}

func validateParam(name string, value *string) bool {
	if value == nil || len(*value) == 0 {
		log.Println("Missing required value for '" + name + "'")
		return false
	}
	return true
}
