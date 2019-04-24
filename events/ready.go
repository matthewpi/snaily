package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/logger"
)

func ReadyEvent(session *discordgo.Session, event *discordgo.Ready) {
	snaily := bot.GetBot()

	if snaily.Config.Discord.Status.Active {
		var gameType discordgo.GameType
		switch snaily.Config.Discord.Status.Type {
		case "playing":
			gameType = discordgo.GameTypeGame
		case "streaming":
			gameType = discordgo.GameTypeStreaming
		case "watching":
			gameType = discordgo.GameTypeWatching
		case "listening":
			gameType = discordgo.GameTypeListening
		}

		err := session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: "playing",
			Game: &discordgo.Game{
				Name: snaily.Config.Discord.Status.Name,
				Type: gameType,
				URL:  "https://krygon.app",
			},
		})
		if err != nil {
			logger.Errorw("[Discord] Failed to update rich presence.", logger.Err(err))
		}
	}
}
