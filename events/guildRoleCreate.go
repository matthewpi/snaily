package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
)

func GuildRoleCreateEvent(_ *discordgo.Session, role *discordgo.GuildRoleCreate) {
	snaily := bot.GetBot()

	if _, ok := snaily.Config.Discord.Guilds[role.GuildID]; !ok {
		return
	}

	// Log the role creation.
	snaily.SendEmbedMessage(
		snaily.Config.Discord.Guilds[role.GuildID].Channels.Punishments,
		0xF8E71C,
		"Role Created",
		fmt.Sprintf("<@%s>", role.Role.ID),
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    fmt.Sprintf("%s#%s", snaily.User.Username, snaily.User.Discriminator),
			IconURL: snaily.User.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{},
		false,
	)
}
