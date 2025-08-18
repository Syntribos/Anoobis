package anoobis_client

import (
	"errors"
	"fmt"
	"github.com/Syntribos/Anoobis/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
	"time"
)

func readHistoricalMessages(dbInfo *storage.DBInfo, session *discordgo.Session, reportChannel string) {
	var lastCheckedMessage string
	var err error
	log.Println(lastCheckedMessage)
	if lastCheckedMessage, err = dbInfo.GetCurrentReportCursor(reportChannel); err != nil {
		log.Printf("Error getting messages from channel " + reportChannel + ":" + err.Error())
		return
	}

	for {
		messages, err := session.ChannelMessages(reportChannel, 100, "", lastCheckedMessage, "")
		fmt.Println(lastCheckedMessage)

		rle := new(discordgo.RateLimitError)
		if errors.As(err, &rle) {
			time.Sleep(rle.RetryAfter)
			continue
		}
		if err != nil {
			log.Printf("An error occurred while trying to retrieve message: %v\n", err)
			return
		}

		if len(messages) == 0 {
			// figure out how to permanently disable this func
			return
		}

		for i := range messages {
			currMessage := messages[len(messages)-i-1].Content
			log.Printf("Retrieved Message Content: %s\n", messages[len(messages)-i-1].Content)
			reg, err := regexp.Compile(`(.*?)(\d{17,18})(.*)`)
			if err != nil {
				log.Printf("An error occurred while trying to compile regex: %v\n", err)
			}
			res := reg.FindStringSubmatch(currMessage)

			for idx, val := range res {
				log.Println("Decomposed: " + strconv.Itoa(idx) + ": " + val)
			}

			lastCheckedMessage = messages[len(messages)-i-1].ID
		}
	}
}
