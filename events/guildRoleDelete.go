package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
)

func GuildRoleDeleteEvent(_ *discordgo.Session, role *discordgo.GuildRoleDelete) {
	snaily := bot.GetBot()

	if _, ok := snaily.Config.Discord.Guilds[role.GuildID]; !ok {
		return
	}

	// Log the role creation.
	snaily.SendEmbedMessage(
		snaily.Config.Discord.Guilds[role.GuildID].Channels.Punishments,
		0xF8E71C,
		"Role Delete",
		fmt.Sprintf("<@%s>", role.RoleID),
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    fmt.Sprintf("%s#%s", snaily.User.Username, snaily.User.Discriminator),
			IconURL: snaily.User.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{},
		false,
	)
}
