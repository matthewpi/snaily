package events

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/logger"
)

func MessageDeleteEvent(_ *discordgo.Session, msg *discordgo.MessageDelete) {
	snaily := bot.GetBot()

	if _, ok := snaily.Config.Discord.Guilds[msg.GuildID]; !ok {
		return
	}

	exec := snaily.Redis.Client.Get(fmt.Sprintf("snaily:message:%s", msg.ID))
	result, err := exec.Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			return
		}

		logger.Errorw("[Discord] Failed to get redis value.", logger.Err(err))
		return
	}

	var originalMessage *discordgo.Message
	if err := json.Unmarshal(result, &originalMessage); err != nil {
		logger.Errorw("[Discord] Failed to json#Unmarshal message.", logger.Err(err))
		return
	}

	// Ignore embeds.
	if len(originalMessage.Embeds) > 0 {
		return
	}

	// Ignore command messages.
	if string(originalMessage.Content[0]) == snaily.Config.Discord.Prefix {
		return
	}

	// Log the message delete.
	snaily.SendEmbedMessage(
		snaily.Config.Discord.Guilds[msg.GuildID].Channels.Messages,
		0xB92222,
		"Message Deleted",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    fmt.Sprintf("%s#%s", originalMessage.Author.Username, originalMessage.Author.Discriminator),
			IconURL: originalMessage.Author.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   "Message ID",
				Value:  msg.ID,
				Inline: false,
			},

			{
				Name:   "Channel",
				Value:  fmt.Sprintf("<#%s> (%s)", originalMessage.ChannelID, originalMessage.ChannelID),
				Inline: false,
			},

			{
				Name:   "Content",
				Value:  originalMessage.Content,
				Inline: false,
			},
		},
		false,
	)
}
