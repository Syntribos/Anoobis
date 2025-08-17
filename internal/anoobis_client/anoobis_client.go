package anoobis_client

import (
	"fmt"
	"github.com/Syntribos/Anoobis/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

type AnoobisClient struct {
}

func Run(botToken *string, guildId *string, dbPath *string) {
	session, _ := discordgo.New("Bot " + *botToken)
	session.Identify.Intents |= discordgo.IntentsAllWithoutPrivileged

	db, err := storage.Init(*dbPath)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("database version: ", db.GetVersionString())

	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer func() {
		err = session.Close()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
