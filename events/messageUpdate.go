package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/logger"
)

func MessageUpdateEvent(session *discordgo.Session, msg *discordgo.MessageUpdate) {
	logger.Infof("[Discord] Message updated: %s", msg.ID)

	if len(msg.EditedTimestamp) > 0 {
		logger.Infof("[Discord] Message modified: %s", msg.ID)
	}
}
