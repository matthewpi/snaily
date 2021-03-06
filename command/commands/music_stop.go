package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
)

// Stop .
func Stop() *command.Command {
	cmd := &command.Command{
		Name:      "stop",
		Aliases:   []string{},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      "boombox",
		Handler:   stopCommandHandler,
	}
	return cmd
}

func stopCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Check if there is not an active stream.
	if snaily.MusicStream[cmd.Message.GuildID] == nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, no music is currently playing.", cmd.Message.Author.ID)
		return
	}

	// Clear the song queue.
	snaily.Queue[cmd.Message.GuildID] = []*bot.Request{}

	// End the music stream.
	snaily.MusicStream[cmd.Message.GuildID].Stop()

	// Send a message to the user.
	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, I've stopped the music.", cmd.Message.Author.ID)
}
