package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/chonla/roman-number-go"
	"github.com/matthewpi/snaily/api/hypixel"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/logger"
	"time"
)

// Hypixel .
func Hypixel() *command.Command {
	cmd := &command.Command{
		Name:    "hypixel",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "id",
				Name:     "id",
				Required: true,
			},
		},
		Enhanced: true,
		Role:     "",
		Handler:  hypixelCommandHandler,
	}
	return cmd
}

func hypixelCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	hypixelApi := &hypixel.API{
		Key: snaily.Config.API.Hypixel.Key,
	}

	player, err := hypixelApi.Player(cmd.Arguments[0])
	if err != nil {
		logger.Errorw("[Hypixel] An error occurred while contacting the api.", logger.Err(err))
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while contacting the Hypixel API.", cmd.Message.Author.ID)
		return
	}

	cmd.SendEmbedMessage(
		cmd.Message.ChannelID,
		0xFFFFFF,
		"",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    player.Displayname,
			IconURL: fmt.Sprintf("https://crafatar.com/avatars/%s", player.UUID),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  player.UUID,
				Inline: false,
			},

			{
				Name:   "Username",
				Value:  player.Displayname,
				Inline: false,
			},

			{
				Name:   "First Joined",
				Value:  time.Unix(0, player.FirstLogin*int64(time.Millisecond)).Format("2006-01-02 15:04:05"),
				Inline: false,
			},

			{
				Name:   "Last Played",
				Value:  time.Unix(0, player.LastLogin*int64(time.Millisecond)).Format("2006-01-02 15:04:05"),
				Inline: false,
			},

			{
				Name: "Pit",
				Value: "```" +
					fmt.Sprintf(`Level: %s-%d
EXP: %d
Renown: %d
Gold: %.2f
`, roman.NewRoman().ToRoman(len(player.Stats.Pit.Profile.Prestiges)), player.Stats.Pit.Profile.XP, player.Stats.Pit.Profile.XP, player.Stats.Pit.Profile.Renown, player.Stats.Pit.Profile.Cash) +
					"```",
				Inline: false,
			},
		},
		true,
	)
}
