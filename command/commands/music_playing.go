package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
)

// Playing .
func Playing() *command.Command {
	cmd := &command.Command{
		Name: "playing",
		Aliases: []string{
			"current",
		},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      "boombox",
		Handler:   playingCommandHandler,
	}
	return cmd
}

func playingCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Check if there is not an active stream.
	if snaily.MusicStream[cmd.Message.GuildID] == nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, no music is currently playing.", cmd.Message.Author.ID)
		return
	}

	// Send a message to the user.
	cmd.SendEmbedMessage(
		cmd.Message.ChannelID,
		0x007EFC,
		"Currently Playing",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    snaily.Config.Build.Name,
			IconURL: snaily.User.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   snaily.Playing[cmd.Message.GuildID].VideoInfo.Title,
				Value:  fmt.Sprintf("Requested By: %s#%s", snaily.Playing[cmd.Message.GuildID].Author.User.Username, snaily.Playing[cmd.Message.GuildID].Author.User.Discriminator),
				Inline: false,
			},
		},
		true,
	)
}
