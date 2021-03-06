package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"
)

// Request represents a music request.
type Request struct {
	Author    *discordgo.Member `json:"author"`
	ChannelID string            `json:"channelId"`
	Video     string            `json:"video"`
	VideoInfo *ytdl.VideoInfo   `json:"videoInfo"`
}

