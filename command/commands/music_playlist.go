package commands

import (
	"fmt"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/logger"
	"github.com/rylio/ytdl"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/url"
)

const (
	// Max number of requests that can be added to the queue per playlist.
	MaxPlaylistItems = 50
	// Max number of search results.
	MaxSearchItems = 5
)

// Playlist .
func Playlist() *command.Command {
	cmd := &command.Command{
		Name:    "playlist",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "youtubePlaylistUrl",
				Name:     "youtube playlist url",
				Required: true,
			},
		},
		Enhanced: false,
		Role:     "boombox",
		Handler:  playlistCommandHandler,
	}
	return cmd
}

func playlistCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	playlistUrl, err := url.Parse(cmd.Arguments[0])
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, invalid playlist url specified.", cmd.Message.Author.ID)
		return
	}

	if len(playlistUrl.Query().Get("list")) < 1 {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, youtube link is not for a playlist.", cmd.Message.Author.ID)
		return
	}

	//client := &http.Client{Transport: &transport.APIKey{Key: snaily.Config.API.Youtube.Key}}
	//service, err := youtube.New(client)
	service, err := youtube.NewService(nil, option.WithAPIKey(snaily.Config.API.Youtube.Key))
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while contacting the Youtube API.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to create Youtube service.", logger.Err(err))
		return
	}

	call := service.PlaylistItems.List("id,snippet").PlaylistId(playlistUrl.Query().Get("list")).MaxResults(MaxPlaylistItems)
	_, err = call.Do()
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred while contacting the Youtube API.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to search for videos.", logger.Err(err))
		return
	}

	results := 0
	err = call.Pages(nil, func(response *youtube.PlaylistItemListResponse) error {
		for _, vid := range response.Items {
			videoInfo, err := ytdl.GetVideoInfo(fmt.Sprintf("https://youtube.com/watch?v=%s", vid.Snippet.ResourceId.VideoId))
			if err != nil {
				continue
			}

			// Check the available video formats.
			formats := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)
			if len(formats) < 1 {
				continue
			}

			// Get the video's download url.
			downloadURL, err := videoInfo.GetDownloadURL(formats[0])
			if err != nil {
				continue
			}

			request := &bot.Request{
				Author:    cmd.Member,
				ChannelID: cmd.Message.ChannelID,
				Video:     downloadURL.String(),
				VideoInfo: videoInfo,
			}
			snaily.AddQueue(request)
			results++

			if results >= MaxPlaylistItems {
				break
			}
		}

		return nil
	})

	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, added %d songs to the queue. (Note: there is a limit of %d songs when adding from a playlist)", cmd.Message.Author.ID, results, MaxPlaylistItems)
}
