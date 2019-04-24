package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"strconv"
)

// Purge .
func Purge() *command.Command {
	cmd := &command.Command{
		Name:    "purge",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "messages",
				Name:     "messages",
				Required: true,
			},
		},
		Enhanced: true,
		Role:     "",
		Handler:  purgeCommandHandler,
	}
	return cmd
}

func purgeCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Convert the argument to an integer.
	messageCount, err := strconv.Atoi(cmd.Arguments[0])
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, please enter a valid integer between 2-100.", cmd.Message.Author.ID)
		return
	}

	// Check if the argument is too low or too high.
	if messageCount < 2 || messageCount > 100 {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, please enter a valid integer between 2-100.", cmd.Message.Author.ID)
		return
	}

	// Fetch a list of message IDs to delete.
	messages, err := snaily.Session.ChannelMessages(cmd.Message.ChannelID, messageCount, "", "", "")
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while fetching a list of messages.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to fetch a list of messages.", logger.Err(err))
		return
	}

	// Convert the []*Message to []string
	var messagesToDelete []string
	for _, msg := range messages {
		messagesToDelete = append(messagesToDelete, msg.ID)
	}

	// Delete the messages.
	err = snaily.Session.ChannelMessagesBulkDelete(cmd.Message.ChannelID, messagesToDelete)
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while deleting messages.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to bulk delete messages.", logger.Err(err))
		return
	}

	// Send a response to the user.
	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, deleted the previous %d messages.", cmd.Message.Author.ID, messageCount)

	// Log the purge.
	cmd.SendEmbedMessage(
		config.Get().Discord.Channels.Messages,
		0xF8E71C,
		"Messages Purge",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    fmt.Sprintf("%s#%s", cmd.Message.Author.Username, cmd.Message.Author.Discriminator),
			IconURL: cmd.Message.Author.AvatarURL(""),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name: "Channel ID",
				Value: cmd.Message.ChannelID,
				Inline: false,
			},
			{
				Name:   "Message Count",
				Value:  cmd.Arguments[0],
				Inline: false,
			},
		},
		false,
	)
}
