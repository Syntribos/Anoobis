package anoobis_client

import (
	"fmt"
	"github.com/Syntribos/Anoobis/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type AnoobisClient struct {
}

func Run(botToken *string, guildId *string, dbPath *string) error {
	fmt.Println("Creating bot session...")
	session, _ := discordgo.New("Bot " + *botToken)
	session.Identify.Intents |= discordgo.IntentsAllWithoutPrivileged
	fmt.Println("Session created.")

	fmt.Println("Initializing database connection...")

	var db *storage.DBInfo
	var err error
	if db, err = storage.Init(*dbPath); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connection established. DB version: ", db.GetVersionString())

	fmt.Println("Connecting to Discord...")

	if err = session.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer func() {
		fmt.Println("Closing Discord session.")
		err = session.Close()
	}()

	fmt.Println("Discord session established. Press CTRL+C to exit.")
	stop := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")
	return err
}
