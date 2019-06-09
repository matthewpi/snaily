package commands

import (
	"fmt"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/logger"
	"github.com/rylio/ytdl"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"strings"
)

// Play .
func Play() *command.Command {
	cmd := &command.Command{
		Name:    "play",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "youtubeUrl",
				Name:     "youtube url or search query",
				Required: true,
			},
		},
		Enhanced: false,
		Role:     "boombox",
		Handler:  playCommandHandler,
	}
	return cmd
}

func playCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	var videoUrl string
	if !strings.HasPrefix(cmd.Argument, "http") {
		service, err := youtube.NewService(nil, option.WithAPIKey(snaily.Config.API.Youtube.Key))
		if err != nil {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while contacting the Youtube API.", cmd.Message.Author.ID)
			logger.Errorw("[Discord] Failed to create Youtube service.", logger.Err(err))
			return
		}

		call := service.Search.List("id,snippet").Q(cmd.Argument).MaxResults(MaxSearchItems)
		_, err = call.Do()
		if err != nil {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while contacting the Youtube API.", cmd.Message.Author.ID)
			logger.Errorw("[Discord] Failed to search for videos.", logger.Err(err))
			return
		}

		err = call.Pages(nil, func(response *youtube.SearchListResponse) error {
			// TODO: Select best video, this might require asking the user, checking views, checking likes, etc.
			for _, vid := range response.Items {
				videoUrl = fmt.Sprintf("https://youtube.com/watch?v=%s", vid.Id.VideoId)
				break
			}

			return nil
		})
		if err != nil {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while finding a video.", cmd.Message.Author.ID)
			logger.Errorw("[Discord] Failed to locate video.", logger.Err(err))
			return
		}
	} else {
		videoUrl = cmd.Arguments[0]
	}

	// Get the video's information.
	videoInfo, err := ytdl.GetVideoInfo(videoUrl)
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
		Author:    cmd.Member,
		ChannelID: cmd.Message.ChannelID,
		Video:     downloadURL.String(),
		VideoInfo: videoInfo,
	}
	snaily.AddQueue(request)

	if snaily.Playing[cmd.Message.GuildID] != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, added \"%s\" to the queue.", cmd.Message.Author.ID, videoInfo.Title)
	}
}
