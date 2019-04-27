package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
)

// Stop .
func Stop() *command.Command {
	cmd := &command.Command{
		Name:      "stop",
		Aliases:   []string{},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      config.Get().Discord.Roles.Boombox,
		Handler:   stopCommandHandler,
	}
	return cmd
}

func stopCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Check if there is not an active stream.
	if snaily.MusicStream == nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, no music is currently playing.", cmd.Message.Author.ID)
		return
	}

	// End the music stream.
	snaily.MusicStream.Stop()

	// TODO: Clear music queue, pause queue from processing?

	// Send a message to the user.
	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, I've stopped the music.", cmd.Message.Author.ID)
}

