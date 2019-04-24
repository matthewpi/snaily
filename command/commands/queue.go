package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
)

// Queue .
func Queue() *command.Command {
	cmd := &command.Command{
		Name:      "queue",
		Aliases:   []string{},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      config.Get().Discord.Roles.Boombox,
		Handler:   queueCommandHandler,
	}
	return cmd
}

func queueCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	if len(cmd.Arguments) < 1 {
		if len(snaily.Queue) < 1 {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, there are currently no songs in the queue.", cmd.Message.Author.ID)
			return
		}

		var fields []*discordgo.MessageEmbedField

		for _, request := range snaily.Queue {
			fields = append(fields,
				&discordgo.MessageEmbedField{
					Name:   request.VideoInfo.Title,
					Value:  fmt.Sprintf("Requested By: %s#%s", request.Author.User.Username, request.Author.User.Discriminator),
					Inline: false,
				},
			)
		}

		cmd.SendEmbedMessage(
			cmd.Message.ChannelID,
			0x007EFC,
			"Music Queue",
			"",
			&discordgo.MessageEmbedAuthor{
				URL:     "",
				Name:    snaily.Config.Build.Name,
				IconURL: snaily.User.AvatarURL(""),
			},
			fields,
			true,
		)
		return
	}
}

