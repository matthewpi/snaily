package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/logger"
)

func MessageUpdateEvent(_ *discordgo.Session, msg *discordgo.MessageUpdate) {
	snaily := bot.GetBot()

	if _, ok := snaily.Config.Discord.Guilds[msg.GuildID]; !ok {
		return
	}

	if len(msg.EditedTimestamp) < 1 {
		return
	}

	logger.Infof("[Discord] Message updated: %s", msg.ID)
}
