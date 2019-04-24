package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/matthewpi/snaily/backend"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

// Bot .
type Bot struct {
	Config      *config.Config       `json:"config"`
	Commands    []*command.Command   `json:"commands"`
	Mongo       *backend.MongoDriver `json:"-"`
	Redis       *backend.RedisDriver `json:"-"`
	Session     *discordgo.Session   `json:"-"`
	User        *discordgo.User      `json:"-"`
	GuildID     string               `json:"guildId"`
	Queue       []*Request           `json:"queue"`
	MusicStream *dca.StreamingSession `json:"-"`
}

// SendMessage .
func (bot *Bot) SendMessage(channelId string, template string, data ...interface{}) *discordgo.Message {
	message := fmt.Sprintf(template, data...)
	msg, err := bot.Session.ChannelMessageSend(channelId, message)
	if err != nil {
		logger.Errorw("[Discord] Failed to send message.", logger.Err(err))
		return nil
	}

	go func() {
		// Wait 10 seconds, then delete the message
		time.Sleep(10 * time.Second)
		bot.DeleteMessage(msg)
	}()

	return msg
}

// SendEmbedMessage .
func (bot *Bot) SendEmbedMessage(channelId string, colour int, title string, description string, author *discordgo.MessageEmbedAuthor, fields []*discordgo.MessageEmbedField, delete bool) *discordgo.Message {
	embed := &discordgo.MessageEmbed{
		Color:       colour,
		Title:       title,
		Description: description,
		Author:      author,
		Fields:      fields,

		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("%s/%s", bot.Config.Build.Commit, bot.Config.Build.Branch),
			IconURL: bot.User.AvatarURL(""),
		},

		Timestamp: time.Now().Format(time.RFC3339),
	}

	msg, err := bot.Session.ChannelMessageSendEmbed(channelId, embed)
	if err != nil {
		logger.Errorw("[Discord] Failed to send embed message.", logger.Err(err))
		return nil
	}

	if delete {
		go func() {
			// Wait 10 seconds, then delete the message
			time.Sleep(10 * time.Second)
			bot.DeleteMessage(msg)
		}()
	}

	return msg
}

// DeleteMessage .
func (bot *Bot) DeleteMessage(message *discordgo.Message) {
	err := bot.Session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "HTTP 404 Not Found") {
			logger.Errorw("[Discord] Failed to delete message.", logger.Err(err))
		}
	}
}

// HasPermission .
func (bot *Bot) HasPermission(member *discordgo.Member, permission int) bool {
	guild, err := snaily.Guild(member.GuildID)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch guild information.", logger.Err(err))
		return false
	}

	if guild.OwnerID == member.User.ID {
		return true
	}

	for _, role := range guild.Roles {
		if !bot.HasRole(member, role.ID) {
			continue
		}

		if role.Permissions&permission == permission {
			return true
		}
	}

	return false
}

// HasRole .
func (bot *Bot) HasRole(member *discordgo.Member, role string) bool {
	for _, roleId := range member.Roles {
		if roleId == role {
			return true
		}
	}

	return false
}

// CanTarget .
func (bot *Bot) CanTarget(member *discordgo.Member, target *discordgo.Member) bool {
	if member.User.ID == target.User.ID {
		return true
	}

	guild, err := snaily.Guild(member.GuildID)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch guild information.", logger.Err(err))
		return false
	}

	if guild.OwnerID == member.User.ID {
		return true
	}

	roles := map[string]*discordgo.Role{}
	for _, role := range guild.Roles {
		roles[role.ID] = role
	}

	var memberRole *discordgo.Role
	for _, roleId := range member.Roles {
		role := roles[roleId]
		if memberRole == nil {
			memberRole = role
		}

		if role.Position > memberRole.Position {
			memberRole = role
		}
	}

	var targetRole *discordgo.Role
	for _, roleId := range target.Roles {
		role := roles[roleId]
		if targetRole == nil {
			targetRole = role
		}

		if role.Position > targetRole.Position {
			targetRole = role
		}
	}

	if targetRole.Position >= memberRole.Position {
		return false
	}

	return true
}

// CanBotTarget .
func (bot *Bot) CanBotTarget(target *discordgo.Member) bool {
	member, err := bot.GuildMember(target.GuildID, bot.User.ID)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch bot member information.", logger.Err(err))
		return false
	}

	return bot.CanTarget(member, target)
}

// Guild .
func (bot *Bot) Guild(guildId string) (*discordgo.Guild, error) {
	// Get the cached guild object.
	guild, err := bot.Session.State.Guild(guildId)
	if err == nil {
		// Return the cached guild object if it exists.
		return guild, nil
	}

	// Get the guild object.
	guild, err = bot.Session.Guild(guildId)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch guild information.", logger.Err(err))
		return nil, err
	}

	// Add the guild to the cache.
	err = bot.Session.State.GuildAdd(guild)
	if err != nil {
		logger.Errorw("[Discord] Failed to add guild to state.", logger.Err(err))
	}

	// Return the guild.
	return guild, err
}

// GuildMember .
func (bot *Bot) GuildMember(guildId string, userId string) (*discordgo.Member, error) {
	// Get the cached member object.
	member, err := bot.Session.State.Member(guildId, userId)
	if err == nil {
		// Return the cached member object if it exists.
		return member, nil
	}

	// Get the guild member object.
	member, err = bot.Session.GuildMember(guildId, userId)
	if err != nil {
		logger.Errorw("[Discord] Failed to fetch guild member information.", logger.Err(err))
		return nil, err
	}

	// Fix because the library doesn't always set this value.
	member.GuildID = guildId

	// Add the member to the cache.
	err = bot.Session.State.MemberAdd(member)
	if err != nil {
		logger.Errorw("[Discord] Failed to add member to state.", logger.Err(err))
	}

	// Return the member.
	return member, err
}

// AddQueue .
func (bot *Bot) AddQueue(request *Request) {
	if bot.Queue == nil {
		bot.Queue = []*Request{}
	}

	// Add the video to the queue.
	bot.Queue = append(bot.Queue, request)
}

func (bot *Bot) Music() {
	go func() {
		for {
			if len(bot.Queue) < 1 {
				time.Sleep(time.Second * 3)
				continue
			}

			var request *Request
			request, bot.Queue = bot.Queue[0], bot.Queue[1:]

			// Get guild information.
			guild, err := snaily.Guild(request.Author.GuildID)
			if err != nil {
				logger.Errorw("[Discord] Failed to get guild information.", logger.Err(err))
				return
			}

			// Loop through connected voice clients.
			var channelId string
			for _, voiceState := range guild.VoiceStates {
				if voiceState.UserID != request.Author.User.ID {
					continue
				}

				// Update the channelId variable.
				channelId = voiceState.ChannelID
			}

			// Check if we found a channel id.
			if len(channelId) == 0 {
				bot.SendMessage(request.ChannelID, "<@%s>, you must be in a voice channel.", request.Author.User.ID)
				return
			}

			// Join voice channel.
			conn, err := snaily.Session.ChannelVoiceJoin(request.Author.GuildID, channelId, false, false)
			if err != nil {
				bot.SendMessage(request.ChannelID, "<@%s>, I cannot join your voice channel.", request.Author.User.ID)
				logger.Errorw("[Discord] Failed to connect to voice channel.", logger.Err(err))
				return
			}

			// Start speaking.
			err = conn.Speaking(true)
			if err != nil {
				logger.Errorw("[Discord] Failed to start speaking.", logger.Err(err))
				return
			}

			// Get the video.
			resp, err := http.Get(request.Video)
			if err != nil {
				bot.SendMessage(request.ChannelID, "An error occurred while downloading the video.")
				logger.Errorw("[Discord] Failed to fetch video.", logger.Err(err))
				return
			}

			options := dca.StdEncodeOptions
			options.RawOutput = true
			options.Bitrate = 96
			options.Application = "lowdelay"
			options.BufferedFrames = 4096

			// Start encoding the video.
			encodingSession, err := dca.EncodeMem(resp.Body, options)
			if err != nil {
				bot.SendMessage(request.ChannelID, "An error occurred while encoding the video.")
				logger.Errorw("[Discord] Failed to encode video.", logger.Err(err))
				return
			}

			// Send a response.
			bot.SendMessage(request.ChannelID, "Now playing \"%s\".", request.VideoInfo.Title)

			// Play the video
			done := make(chan error)
			stream := dca.NewStream(encodingSession, conn, done)
			bot.MusicStream = stream
			err = <-done
			if err != nil && err != io.EOF {
				bot.SendMessage(request.ChannelID, "An error occurred during playback.")
				logger.Errorw("[Discord] Error occurred while playing video.", logger.Err(err))
				return
			}

			// Stop speaking.
			err = conn.Speaking(false)
			if err != nil {
				logger.Errorw("[Discord] Failed to stop speaking.", logger.Err(err))
				return
			}

			_, err = stream.Finished()
			if err != nil {
				logger.Errorw("[Discord] Failed to close stream.", logger.Err(err))
			}
			bot.MusicStream = nil

			encodingSession.Stop()
			encodingSession.Cleanup()

			if len(bot.Queue) < 1 {
				err = conn.Disconnect()
				if err != nil {
					logger.Errorw("[Discord] Failed to disconnect from voice channel.", logger.Err(err))
					return
				}
			}
		}
	}()
}

var snaily *Bot

// GetBot .
func GetBot() *Bot {
	return snaily
}

// SetBot .
func SetBot(b *Bot) {
	snaily = b
}
