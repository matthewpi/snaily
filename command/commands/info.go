package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"runtime"
	"strconv"
)

// Info .
func Info() *command.Command {
	cmd := &command.Command{
		Name:      "info",
		Aliases:   []string{"status"},
		Arguments: []*command.Argument{},
		Enhanced:  false,
		Role:      "",
		Handler:   infoCommandHandler,
	}
	return cmd
}

func infoCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	stats, _ := snaily.Session.State.Guild(cmd.Message.GuildID)

	guildCount := strconv.Itoa(len(snaily.Session.State.Guilds))
	memberCount := strconv.Itoa(len(stats.Members))
	channelCount := strconv.Itoa(len(stats.Channels))

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	systemInfo := fmt.Sprintf(
		"Go Version:   %s\nGo Routines:  %d\nMemory Usage: %s / %s\nGarbage:      %s",
		snaily.Config.Build.GoVersion,
		runtime.NumGoroutine(),
		fmt.Sprintf("%d MB", memStats.Alloc/1024/1024),
		fmt.Sprintf("%d MB", memStats.TotalAlloc/1024/1024),
		fmt.Sprintf("%d MB", memStats.GCSys/1024/1024),
	)

	cmd.SendEmbedMessage(
		cmd.Message.ChannelID,
		0x007EFC,
		"",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    snaily.Config.Build.Name,
			IconURL: snaily.User.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   "Guild Count",
				Value:  guildCount,
				Inline: true,
			},

			{
				Name:   "Channel Count",
				Value:  channelCount,
				Inline: true,
			},

			{
				Name:   "Member Count",
				Value:  memberCount,
				Inline: false,
			},

			{
				Name:   "System Information",
				Value:  "```" + systemInfo + "```",
				Inline: false,
			},
		},
		true,
	)
}
