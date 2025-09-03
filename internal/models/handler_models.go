package models

import "fmt"

const msgLinkFmt = "https://discord.com/channels/%s/%s/%s"

type UserReport struct {
	UserId    string
	MessageId string
}

type ReportPackage struct {
	Reports   []UserReport
	Reason    string
	GuildId   string
	ChannelId string
}

func (report ReportPackage) GetMessageLink(idIndex int) string {
	return fmt.Sprintf(msgLinkFmt, report.GuildId, report.ChannelId, report.Reports[idIndex].MessageId)
}
