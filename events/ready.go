package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/logger"
)

func ReadyEvent(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "playing",
		Game: &discordgo.Game{
			Name: "with mods",
			Type: discordgo.GameTypeGame,
			URL:  "https://krygon.app",
		},
	})

	if err != nil {
		logger.Errorw("[Discord] Failed to update rich presence.", logger.Err(err))
	}
}
