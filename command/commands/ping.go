package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"time"
)

// Ping .
func Ping() *command.Command {
	cmd := &command.Command{
		Name:      "ping",
		Aliases:   []string{},
		Arguments: []*command.Argument{},
		Enhanced:  true,
		Role:      "",
		Handler:   pingCommandHandler,
	}
	return cmd
}

func pingCommandHandler(cmd *command.Execution) {
	stacktraceBot := bot.GetBot()
	cmd.SendMessage(cmd.Message.ChannelID, "Pong! %vms", int64(stacktraceBot.Session.HeartbeatLatency()/time.Millisecond))
}
