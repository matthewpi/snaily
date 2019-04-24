package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"strconv"
	"strings"
	"time"
)

// Argument represents a command argument.
type Argument struct {
	ID       string
	Name     string
	Required bool
}

// Command represents a command and it's handler.
type Command struct {
	Name      string
	Aliases   []string
	Arguments []*Argument
	Enhanced  bool
	Role      string
	Handler   func(cmd *Execution)
}

// Usage returns a string with the proper command usage.
func (command *Command) Usage() string {
	var usage strings.Builder

	for _, arg := range command.Arguments {
		usage.WriteString(" ")

		if arg.Required {
			usage.WriteString("<")
		} else {
			usage.WriteString("[")
		}

		usage.WriteString(arg.Name)

		if arg.Required {
			usage.WriteString(">")
		} else {
			usage.WriteString("]")
		}
	}

	return usage.String()
}

// Execution represents a command execution.
type Execution struct {
	Label     string
	Argument  string
	Arguments []string
	Args      map[string]string
	Member    *discordgo.Member
	Message   *discordgo.Message
	Session   *discordgo.Session
	BotUser   *discordgo.User
	Command   *Command
	Messages  []string
}

// SendMessage .
func (exec *Execution) SendMessage(channelId string, template string, data ...interface{}) {
	message := fmt.Sprintf(template, data...)
	msg, err := exec.Session.ChannelMessageSend(channelId, message)
	if err != nil {
		logger.Errorw("[Discord] Failed to send message.", logger.Err(err))
		return
	}

	// Add the message to the messages list.
	exec.Messages = append(exec.Messages, msg.ID)
}

// SendEmbedMessage .
func (exec *Execution) SendEmbedMessage(channelId string, colour int, title string, description string, author *discordgo.MessageEmbedAuthor, fields []*discordgo.MessageEmbedField, delete bool) {
	embed := &discordgo.MessageEmbed{
		Color:       colour,
		Title:       title,
		Description: description,
		Author:      author,
		Fields:      fields,

		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("%s/%s", config.Get().Build.Commit, config.Get().Build.Branch),
			IconURL: exec.BotUser.AvatarURL(""),
		},

		Timestamp: time.Now().Format(time.RFC3339),
	}

	msg, err := exec.Session.ChannelMessageSendEmbed(channelId, embed)
	if err != nil {
		logger.Errorw("[Discord] Failed to send embed message.", logger.Err(err))
		return
	}

	if delete {
		// Add the message to the messages list.
		exec.Messages = append(exec.Messages, msg.ID)
	}
}

// DeleteMessage .
func (exec *Execution) DeleteMessage(message *discordgo.Message) {
	err := exec.Session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "HTTP 404 Not Found") {
			logger.Errorw("[Discord] Failed to delete message.", logger.Err(err))
		}
	}
}

// GetMember .
func (exec *Execution) GetMember(targetString string) *discordgo.Member {
	var target *discordgo.Member

	if string(targetString[0]) == "<" && string(targetString[1]) == "@" && string(targetString[len(targetString)-1]) == ">" {
		targetString = targetString[2 : len(targetString)-1]
	} else if _, err := strconv.ParseInt(targetString, 10, 64); err == nil {
		// argument is a user id
	} else {
		guild, err := exec.Guild(exec.Message.GuildID)
		if err != nil {
			exec.SendMessage(exec.Message.ChannelID, "<@%s>, an error occurred while getting guild information.", exec.Message.Author.ID)
			logger.Errorw("[Discord] Failed to fetch guild information.", logger.Err(err))
			return nil
		}

		found := false
		for _, member := range guild.Members {
			if strings.HasPrefix(member.User.Username, targetString) {
				targetString = member.User.ID
				found = true
				break
			}
		}

		if !found {
			exec.SendMessage(exec.Message.ChannelID, "<@%s>, no matching target found.", exec.Message.Author.ID)
			return nil
		}
	}

	target, err := exec.GuildMember(exec.Message.GuildID, targetString)
	if err != nil {
		exec.SendMessage(exec.Message.ChannelID, "<@%s>, an error occurred while getting target information.", exec.Message.Author.ID)
		logger.Errorw("[Discord] Failed to obtain account details for target.", logger.Err(err))
		return nil
	}

	if target == nil {
		exec.SendMessage(exec.Message.ChannelID, "<@%s>, no matching target found.", exec.Message.Author.ID)
		return nil
	}

	return target
}

// Guild attempts to return a guild object from the cache, if the object is not cached it will use the discord api.
func (exec *Execution) Guild(guildId string) (*discordgo.Guild, error) {
	// Get the cached guild object.
	guild, err := exec.Session.State.Guild(guildId)
	if err == nil {
		// Return the cached guild object if it exists.
		return guild, nil
	}

	// Get the guild object.
	guild, err = exec.Session.Guild(guildId)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch guild information.", logger.Err(err))
		return nil, err
	}

	// Add the guild to the cache.
	err = exec.Session.State.GuildAdd(guild)
	if err != nil {
		logger.Errorw("[Discord] Failed to add guild to state.", logger.Err(err))
	}

	// Return the guild.
	return guild, err
}

// GuildMember attempts to return a guild member object from the cache, if the object is not cached it will use the discord api.
func (exec *Execution) GuildMember(guildId string, userId string) (*discordgo.Member, error) {
	// Get the cached member object.
	member, err := exec.Session.State.Member(guildId, userId)
	if err == nil {
		// Return the cached member object if it exists.
		return member, nil
	}

	// Get the guild member object.
	member, err = exec.Session.GuildMember(guildId, userId)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch guild member information.", logger.Err(err))
		return nil, err
	}

	// Fix because the library doesn't always set this value.
	member.GuildID = guildId

	// Add the member to the cache.
	err = exec.Session.State.MemberAdd(member)
	if err != nil {
		logger.Errorw("[Discord] Failed to add member to state.", logger.Err(err))
	}

	// Return the member.
	return member, err
}

// Cleanup cleans up all messages the client and bot have sent during the command execution.
func (exec *Execution) Cleanup() {
	// Append the command message to the message list
	exec.Messages = append([]string{exec.Message.ID}, exec.Messages...)

	go func() {
		// Wait 10 seconds, then delete the messages
		time.Sleep(10 * time.Second)

		err := exec.Session.ChannelMessagesBulkDelete(exec.Message.ChannelID, exec.Messages)
		if err != nil {
			logger.Errorw("[Discord] Failed to bulk delete messages.", logger.Err(err))
		}
	}()
}
