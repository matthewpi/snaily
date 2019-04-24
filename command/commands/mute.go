package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"github.com/matthewpi/snaily/utils"
	"strings"
	"time"
)

// Mute .
func Mute() *command.Command {
	cmd := &command.Command{
		Name:    "mute",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "target",
				Name:     "target",
				Required: true,
			},
			{
				ID:       "duration",
				Name:     "duration",
				Required: true,
			},
			{
				ID:       "reason",
				Name:     "reason",
				Required: true,
			},
		},
		Enhanced: true,
		Role:     "",
		Handler:  muteCommandHandler,
	}
	return cmd
}

func muteCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	target := cmd.GetMember(cmd.Arguments[0])

	// Check if the target is nil, the function already handles responses.
	if target == nil {
		return
	}

	// Check if the command sender can target the selected user
	if !snaily.CanTarget(cmd.Member, target) {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, you may not target this user.", cmd.Message.Author.ID)
		return
	}

	// Check if the user is targeting themselves or the bot user.
	if target.User.ID == cmd.Message.Author.ID || target.User.ID == snaily.User.ID {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, you might as well start a war.", cmd.Message.Author.ID)
		return
	}

	// Check if the bot can target the selected user
	if !snaily.CanBotTarget(target) {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, I cannot target that user.", cmd.Message.Author.ID)
		return
	}

	var builder strings.Builder
	for index, arg := range cmd.Arguments[2:] {
		if index > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(arg)
	}
	reason := builder.String()

	var duration time.Duration
	var err error
	if cmd.Arguments[1] == "0" {
		duration = 0
	} else {
		duration, err = utils.ParseDuration(cmd.Arguments[1])
		if err != nil {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, please enter a valid duration.", cmd.Message.Author.ID)
			return
		}
	}

	err = snaily.Session.GuildMemberRoleAdd(cmd.Message.GuildID, target.User.ID, config.Get().Discord.Roles.Muted)
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while muting the selected user.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to add role to user.", logger.Err(err))
		return
	}

	// Log the mute.
	cmd.SendEmbedMessage(
		config.Get().Discord.Channels.Punishments,
		0xF8E71C,
		"Mute",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    fmt.Sprintf("%s#%s", cmd.Message.Author.Username, cmd.Message.Author.Discriminator),
			IconURL: cmd.Message.Author.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   "Punisher",
				Value:  fmt.Sprintf("%s#%s (%s)", cmd.Message.Author.Username, cmd.Message.Author.Discriminator, cmd.Message.Author.ID),
				Inline: false,
			},

			{
				Name:   "Target",
				Value:  fmt.Sprintf("%s#%s (%s)", target.User.Username, target.User.Discriminator, target.User.ID),
				Inline: false,
			},

			{
				Name:   "Reason",
				Value:  reason,
				Inline: false,
			},

			{
				Name:   "Duration",
				Value:  utils.DurationString(duration),
				Inline: false,
			},
		},
		false,
	)

	// Send a response.
	cmd.SendMessage(cmd.Message.ChannelID,
		"<@%s>, %s#%s (%s) has been muted for \"%s\".", cmd.Message.Author.ID, target.User.Username, target.User.Discriminator, target.User.ID, reason,
	)
}
