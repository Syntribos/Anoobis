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

func Run(botToken *string, guildId *string, channelId *string, dbPath *string) error {
	log.Println("Creating bot session...")
	session, _ := discordgo.New("Bot " + *botToken)
	session.Identify.Intents |= discordgo.IntentsAllWithoutPrivileged
	session.Identify.Intents |= discordgo.PermissionAddReactions
	session.Identify.Intents |= discordgo.PermissionEmbedLinks
	session.Identify.Intents |= discordgo.PermissionReadMessageHistory
	session.Identify.Intents |= discordgo.PermissionSendMessages

	//Not ok for some reason?
	//session.Identify.Intents |= discordgo.PermissionAttachFiles
	//session.Identify.Intents |= discordgo.PermissionSendMessagesInThreads
	log.Println("Session created.")

	log.Println("Initializing database connection...")

	var db *storage.DBInfo
	var err error
	if db, err = storage.Init(*dbPath); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established. DB version: ", db.GetVersionString())

	defer func() {
		log.Println("Closing Discord session.")
		err = session.Close()
	}()

	session.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		readHistoricalMessages(db, session, *channelId)
	})

	log.Println("Connecting to Discord...")
	if err = session.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	fmt.Println("Discord session established. Press CTRL+C to exit.")
	stop := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")
	return err
}
