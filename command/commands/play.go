package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"github.com/rylio/ytdl"
)

// Play .
func Play() *command.Command {
	cmd := &command.Command{
		Name:    "play",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "youtubeUrl",
				Name:     "youtube url",
				Required: true,
			},
		},
		Enhanced: false,
		Role:     config.Get().Discord.Roles.Boombox,
		Handler:  playCommandHandler,
	}
	return cmd
}

func playCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	// Get the video's information.
	videoInfo, err := ytdl.GetVideoInfo(cmd.Arguments[0])
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while loading video information.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to load youtube video.", logger.Err(err))
		return
	}

	// Check the available video formats.
	formats := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)

	if len(formats) < 1 {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, invalid video url.", cmd.Message.Author.ID)
		return
	}

	// Get the video's download url.
	downloadURL, err := videoInfo.GetDownloadURL(formats[0])
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while getting the video url.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to get video download url.", logger.Err(err))
		return
	}

	request := &bot.Request{
		Author: cmd.Member,
		ChannelID: cmd.Message.ChannelID,
		Video: downloadURL.String(),
		VideoInfo: videoInfo,
	}
	snaily.AddQueue(request)

	if len(snaily.Queue) > 2 {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, added \"%s\" to the queue.", cmd.Message.Author.ID, videoInfo.Title)
	}
}
