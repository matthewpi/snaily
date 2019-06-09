package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
)

// Shuffle .
func Shuffle() *command.Command {
	cmd := &command.Command{
		Name:      "shuffle",
		Aliases:   []string{},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      "boombox",
		Handler:   shuffleCommandHandler,
	}
	return cmd
}

func shuffleCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Toggle the shuffle state.
	snaily.Shuffle[cmd.Message.GuildID] = !snaily.Shuffle[cmd.Message.GuildID]

	var message string
	if snaily.Shuffle[cmd.Message.GuildID] {
		message = "enabled"
	} else {
		message = "disabled"
	}

	// Send a message to the user.
	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, shuffle %s.", cmd.Message.Author.ID, message)
}
