package events

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"strings"
	"time"
	"unicode"
)

func MessageCreateEvent(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// Get the stored stacktrace.bot object.
	snaily := bot.GetBot()

	// Check if the message author is a bot.
	if msg.Author.Bot {
		return
	}

	// Get the configured command prefix.
	prefix := config.Get().Discord.Prefix

	// Check if the message starts with the configured command prefix.
	if string(msg.Content[0]) != prefix {
		// Beware, profane language D:
		// Also I do realize this is the dumbest way to check for profanity
		// and that I could make this config based, but this is for testing.
		// there is also like 8 million ways to bypass this, but it will prevent
		// most people from being able to do so.
		if strings.Contains(msg.Content, "nigg") || strings.Contains(msg.Content, "fagg") || strings.Contains(msg.Content, "spick") || strings.Contains(msg.Content, "beaner") {
			snaily.DeleteMessage(msg.Message)

			// Log the message delete.
			snaily.SendEmbedMessage(
				config.Get().Discord.Channels.Messages,
				0xB92222,
				"Message Deleted",
				"",
				&discordgo.MessageEmbedAuthor{
					URL:     "",
					Name:    fmt.Sprintf("%s#%s", msg.Author.Username, msg.Author.Discriminator),
					IconURL: msg.Author.AvatarURL(""),
				},
				[]*discordgo.MessageEmbedField{
					{
						Name:   "Message ID",
						Value:  msg.ID,
						Inline: false,
					},

					{
						Name:   "Channel",
						Value:  fmt.Sprintf("<#%s>", msg.ChannelID),
						Inline: false,
					},

					{
						Name:   "Content",
						Value:  msg.Content,
						Inline: false,
					},

					{
						Name:   "Reason",
						Value:  "Profane Language",
						Inline: false,
					},
				},
				false,
			)
			return
		}

		go func() {
			messageJson, err := json.Marshal(msg.Message)
			if err != nil {
				logger.Errorw("[Discord] Failed to json#Marshal message.", logger.Err(err))
				return
			}

			exec := snaily.Redis.Client.Set(fmt.Sprintf("snaily:message:%s", msg.ID), messageJson, 0)
			if exec.Err() != nil {
				logger.Errorw("[Discord] Failed to set redis value.", logger.Err(err))
				return
			}
		}()
		return
	}

	// Get the message's content.
	content := msg.Content

	var label string
	var argument string
	for index, char := range msg.Content {
		// Check if this is the last character in the array.
		if index+1 == len(content) {
			label = strings.ToLower(content[1 : index+1])
			break
		}

		// Check if the character is a space.
		if unicode.IsSpace(char) {
			label = strings.ToLower(content[1:index])
			argument = content[index+1:]
			break
		}
	}

	// Get the matching command object using the user's message.
	var cmd *command.Command
	for _, c := range snaily.Commands {
		if label == c.Name {
			cmd = c
			break
		}

		for _, alias := range c.Aliases {
			if label == alias {
				cmd = c
				break
			}
		}

		if cmd != nil {
			break
		}
	}

	// Check if the found command object is nil.
	if cmd == nil {
		return
	}

	// Handle parsing command arguments.
	var arguments []string
	if strings.Contains(argument, " ") {
		arguments = strings.Split(argument, " ")
	} else if len(argument) > 0 {
		arguments = []string{argument}
	} else {
		arguments = []string{}
	}

	// Get the user as a guild member.
	member, err := snaily.GuildMember(msg.GuildID, msg.Author.ID)
	if err != nil {
		snaily.SendMessage(msg.ChannelID, "<@%s>, an error occurred while loading your user information, I guess you don't exist.", msg.Author.ID)
		return
	}

	// Check if the user does not have administrator permissions.
	if !snaily.HasPermission(member, discordgo.PermissionAdministrator) {
		// Check if the command requires enhanced permissions.
		if cmd.Enhanced {
			// Check if the user does not have the "Enhanced" role.
			if !snaily.HasRole(member, config.Get().Discord.Roles.Enhanced) {
				snaily.SendMessage(msg.ChannelID, "<@%s>, no permission.", msg.Author.ID)
				return
			}
		}

		// Check if the command requires a specific role.
		if len(cmd.Role) > 1 {
			// Check if the user does not have the configured role.
			if !snaily.HasRole(member, cmd.Role) {
				snaily.SendMessage(msg.ChannelID, "<@%s>, no permission.", msg.Author.ID)
				return
			}
		}
	}

	// Get the number of required arguments.
	required := 0
	for _, arg := range cmd.Arguments {
		if !arg.Required {
			continue
		}

		required++
	}

	// Check if the required amount of arguments were not met.
	if len(arguments) < required {
		go func() {
			// Wait 10 seconds, then delete the message
			time.Sleep(10 * time.Second)
			snaily.DeleteMessage(msg.Message)
		}()

		snaily.SendMessage(msg.ChannelID, "<@%s>, usage: `%s%s%s`", msg.Author.ID, prefix, label, cmd.Usage())
		return
	}

	// Create a new command execution object.
	execution := &command.Execution{
		Label:     label,
		Argument:  argument,
		Arguments: arguments,
		Member:    member,
		Message:   msg.Message,
		Session:   snaily.Session,
		BotUser:   snaily.User,
		Command:   cmd,
	}

	// Call the command handler.
	cmd.Handler(execution)

	// Call the execution's cleanup method.
	execution.Cleanup()
}
