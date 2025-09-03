package anoobis_client

import (
	"errors"
	"fmt"
	m "github.com/Syntribos/Anoobis/internal/models"
	"github.com/Syntribos/Anoobis/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
	"time"
)

func readHistoricalMessages(dbInfo *m.DBInfo, session *discordgo.Session, reportGuild string, reportChannel string) {
	var lastCheckedMessage string
	var err error
	log.Println(lastCheckedMessage)
	if lastCheckedMessage, err = storage.GetCurrentReportCursor(dbInfo, reportChannel); err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			log.Printf("Error getting messages from channel " + reportChannel + ":" + err.Error())
			return
		}
		log.Printf("No scan history found, starting from the beginning.")
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

		for i := 0; i < len(messages); i++ {
			currMessage := messages[len(messages)-i-1]
			reportAuthor := currMessage.Author.ID
			reg, err := regexp.Compile(`(.*)?(?:[^\d@&<]|^)(\d{17,18})(?:\D|$)(.*)?`)
			if err != nil {
				log.Printf("An error occurred while trying to compile regex: %v\n", err)
				continue
			}
			res := reg.FindStringSubmatch(currMessage.Content)

			if len(res) == 0 {
				continue
			}

			reports := []m.UserReport{
				{
					UserId:    res[2],
					MessageId: currMessage.ID,
				},
			}

			reason := strings.TrimSpace(strings.TrimSpace(res[1]) + " " + strings.TrimSpace(res[3]))
			if len(reason) == 0 {
				j := i + 1
				nextMessage := messages[len(messages)-j-1]
				for j < len(messages) && nextMessage.Author.ID == reportAuthor {
					res = reg.FindStringSubmatch(nextMessage.Content)

					if len(res) == 0 {
						reason = nextMessage.Content
						break
					}

					newReport := m.UserReport{
						UserId:    res[2],
						MessageId: nextMessage.ID,
					}

					reports = append(reports, newReport)

					j++
					nextMessage = messages[len(messages)-j-1]
				}

				i = j
			}

			reportPackage := m.ReportPackage{
				Reports:   reports,
				Reason:    reason,
				GuildId:   reportGuild,
				ChannelId: currMessage.ChannelID,
			}

			if storage.SaveReport(dbInfo, reportPackage) != nil {
				panic(err)
			}

			lastCheckedMessage = messages[len(messages)-i-1].ID
		}

		err = storage.SaveReportCursor(dbInfo, lastCheckedMessage)
		if err != nil {
			log.Print("An error occurred while trying to save the last read report: " + err.Error())
			return
		}
	}
}
