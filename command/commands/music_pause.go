package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
)

// Pause .
func Pause() *command.Command {
	cmd := &command.Command{
		Name:      "pause",
		Aliases:   []string{},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      "boombox",
		Handler:   pauseCommandHandler,
	}
	return cmd
}

func pauseCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Check if there is not an active stream.
	if snaily.MusicStream[cmd.Message.GuildID] == nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, no music is currently playing.", cmd.Message.Author.ID)
		return
	}

	// Toggle the stream's pause state.
	snaily.MusicStream[cmd.Message.GuildID].SetPaused(!snaily.MusicStream[cmd.Message.GuildID].Paused())

	// Check if the stream is paused.
	if snaily.MusicStream[cmd.Message.GuildID].Paused() {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, I've paused the music.", cmd.Message.Author.ID)
	} else {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, I've resumed the music.", cmd.Message.Author.ID)
	}
}
