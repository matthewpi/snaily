package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
)

// Skip .
func Skip() *command.Command {
	cmd := &command.Command{
		Name: "skip",
		Aliases: []string{
			"next",
		},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      "boombox",
		Handler:   skipCommandHandler,
	}
	return cmd
}

func skipCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Check if there is not an active stream.
	if snaily.MusicStream[cmd.Message.GuildID] == nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, no music is currently playing.", cmd.Message.Author.ID)
		return
	}

	// End the music stream.
	snaily.MusicStream[cmd.Message.GuildID].Stop()

	// Send a message to the user.
	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, song skipped.", cmd.Message.Author.ID)
}
