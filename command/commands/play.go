package commands

import (
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/logger"
	"github.com/matthewpi/snaily/music"
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
	stacktraceBot := bot.GetBot()

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

	request := &music.Request{
		Author: cmd.Member,
		ChannelID: cmd.Message.ChannelID,
		Video: downloadURL.String(),
		VideoInfo: videoInfo,
	}
	stacktraceBot.AddQueue(request)

	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, added \"%s\" to the queue.", cmd.Message.Author.ID, videoInfo.Title)

	/*// Get the video.
	resp, err := http.Get(downloadURL.String())
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, Error: failed to download video.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to fetch video.", logger.Err(err))
		return
	}

	// Start encoding the video.
	encodingSession, err := dca.EncodeMem(resp.Body, options)
	if err != nil {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, Error: failed to encode video.", cmd.Message.Author.ID)
		logger.Errorw("[Discord] Failed to encode video.", logger.Err(err))
		return
	}
	defer encodingSession.Cleanup()

	// Send a response.
	cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, now playing \"%s\"", cmd.Message.Author.ID, videoInfo.Title)

	// Play the video
	done := make(chan error)
	stream := dca.NewStream(encodingSession, conn, done)
	err = <-done
	if err != nil && err != io.EOF {
		cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, an error occurred during playback.", cmd.Message.Author.ID)
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

	encodingSession.Stop()
	encodingSession.Cleanup()

	// Disconnect from voice.
	err = conn.Disconnect()
	if err != nil {
		logger.Errorw("[Discord] Failed to disconnect from voice channel.", logger.Err(err))
		return
	}*/
}
