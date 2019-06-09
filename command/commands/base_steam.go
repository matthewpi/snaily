package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/logger"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// Steam .
func Steam() *command.Command {
	cmd := &command.Command{
		Name:    "steam",
		Aliases: []string{},
		Arguments: []*command.Argument{
			{
				ID:       "id;url",
				Name:     "input",
				Required: true,
			},
		},
		Enhanced: false,
		Role:     "",
		Handler:  steamCommandHandler,
	}
	return cmd
}

func steamCommandHandler(cmd *command.Execution) {
	snaily := bot.GetBot()

	var steam64 string
	vanityUrl := false

	if strings.HasPrefix(cmd.Argument, "http") {
		// Profile URL or Custom URL
		domainEnd := strings.LastIndex(cmd.Argument, ".") + 5

		sections := strings.Split(cmd.Argument[domainEnd:], "/")
		if sections[0] == "id" {
			vanityUrl = true
		}
		steam64 = sections[1]
	} else if strings.HasPrefix(cmd.Argument, "STEAM_") {
		// Steam ID 2
		sections := strings.Split(cmd.Argument[6:], ":")

		if len(sections) != 3 {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, that is not in the proper format.", cmd.Message.Author.ID)
			return
		}

		magic, _ := new(big.Int).SetString("76561197960265728", 10)
		steam, _ := new(big.Int).SetString(sections[2], 10)
		steam = steam.Mul(steam, big.NewInt(2))
		steam = steam.Add(steam, magic)
		auth, _ := new(big.Int).SetString(sections[1], 10)
		steam64 = steam.Add(steam, auth).String()
	} else if _, err := strconv.ParseInt(cmd.Argument, 10, 64); err == nil {
		// Steam 64
		steam64 = cmd.Argument
	} else {
		// Vanity URL or invalid input
		vanityUrl = true
		steam64 = strings.TrimSpace(cmd.Argument)
	}

	// Check if the input was a vanity url.
	if vanityUrl {
		// Contact the steam api to resolve the vanity url.
		vanityResponse, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s&url_type=1", snaily.Config.API.Steam.Key, steam64))
		if err != nil {
			logger.Errorw("[Discord] Failed to load steam user information.", logger.Err(err))
			return
		}

		// Read the response body.
		vanityContent, err := ioutil.ReadAll(vanityResponse.Body)
		if err != nil {
			logger.Errorw("[Discord] Failed to read response body.", logger.Err(err))
			return
		}

		err = vanityResponse.Body.Close()
		if err != nil {
			logger.Errorw("[Discord] Failed to close response body.", logger.Err(err))
			return
		}

		// Decode the response body into a JSON object.
		var vanityData *VanityResponse
		err = json.Unmarshal(vanityContent, &vanityData)
		if err != nil {
			logger.Errorw("[Discord] Failed to unmarshal json response.", logger.Err(err))
			return
		}

		// Check if the response returned a steam id.
		if len(vanityData.Response.SteamID) == 0 || vanityData.Response.Message == "No match" {
			cmd.SendMessage(cmd.Message.ChannelID, "<@%s>, I failed to find an account with that vanity url.", cmd.Message.Author.ID)
			return
		}

		// Update the variable with the returned steam id.
		steam64 = vanityData.Response.SteamID
	}

	// Contact the steam api to get profile information.
	response, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", snaily.Config.API.Steam.Key, steam64))
	if err != nil {
		logger.Errorw("[Discord] Failed to load steam user information.", logger.Err(err))
		return
	}

	// Read the response body.
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Errorw("[Discord] Failed to read response body.", logger.Err(err))
		return
	}

	err = response.Body.Close()
	if err != nil {
		logger.Errorw("[Discord] Failed to close response body.", logger.Err(err))
		return
	}

	// Decode the response body into a JSON object.
	var users *SteamUsers
	err = json.Unmarshal(contents, &users)
	if err != nil {
		logger.Errorw("[Discord] Failed to unmarshal json response.", logger.Err(err))
		return
	}
	user := users.Response.Players[0]

	// Offline
	color := 0x898989

	// Online
	if user.Status == 1 {
		color = 0x57CBDE
	}

	// In-Game
	if len(user.GameID) > 0 {
		color = 0x90BA3C
	}

	cmd.SendEmbedMessage(
		cmd.Message.ChannelID,
		color,
		"",
		"",
		&discordgo.MessageEmbedAuthor{
			URL:     "",
			Name:    user.Name,
			IconURL: user.Avatar,
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  user.Name,
				Inline: false,
			},

			{
				Name:   "Steam ID",
				Value:  user.SteamID(),
				Inline: false,
			},

			{
				Name:   "Steam 3",
				Value:  user.Steam3(),
				Inline: false,
			},

			{
				Name:   "Steam 64",
				Value:  fmt.Sprintf("%s (%s)", user.Steam64, user.Steam64Hex()),
				Inline: false,
			},
		},
		true,
	)
}

type VanityResponse struct {
	Response struct {
		SteamID string `json:"steamid"`
		Success int    `json:"success"`
		Message string `json:"message"`
	} `json:"response"`
}

type SteamUsers struct {
	Response struct {
		Players []SteamUser `json:"players"`
	} `json:"response"`
}

type SteamUser struct {
	Steam64             string `json:"steamid"`
	CommunityVisibility int    `json:"communityvisibilitystate"`
	ProfileState        int    `json:"profilestate"`
	Name                string `json:"personaname"`
	LastLogoff          int64  `json:"lastlogoff"`
	CommentPermission   int    `json:"commentpermission"`
	ProfileURL          string `json:"profileurl"`
	Avatar              string `json:"avatar"`
	AvatarMedium        string `json:"avatarmedium"`
	AvatarFull          string `json:"avatarfull"`
	Status              int    `json:"personastate"`
	RealName            string `json:"realname"`
	PrimaryClanID       string `json:"primaryclanid"`
	CreatedAt           int64  `json:"timecreated"`
	Country             string `json:"loccountrycode"`
	State               string `json:"locstatecode"`
	City                int64  `json:"loccityid"`
	GameID              string `json:"gameid"`
	GameServerIP        string `json:"gameserverip"`
	GameExtraInfo       string `json:"gameextrainfo"`
}

func (user *SteamUser) SteamID() string {
	steam64, err := strconv.ParseInt(user.Steam64, 10, 64)
	if err != nil {
		logger.Errorw("[Discord] Failed to convert Steam 64 to int64.", logger.Err(err))
		return "STEAM_0:0:0"
	}

	steamID := new(big.Int).SetInt64(steam64)
	magic, _ := new(big.Int).SetString("76561197960265728", 10)
	steamID = steamID.Sub(steamID, magic)
	isServer := new(big.Int).And(steamID, big.NewInt(1))
	steamID = steamID.Sub(steamID, isServer)
	steamID = steamID.Div(steamID, big.NewInt(2))
	return "STEAM_0:" + isServer.String() + ":" + steamID.String()
}

func (user *SteamUser) Steam3() string {
	steamId := user.SteamID()

	steamIDParts := strings.Split(string(steamId), ":")
	steamLastPart, err := strconv.ParseUint(string(steamIDParts[len(steamIDParts)-1]), 10, 64)
	if err != nil {
		logger.Errorw("[Discord] Failed to convert Steam ID to uint.", logger.Err(err))
		return "[U:1:0]"
	}

	return "[U:1:" + strconv.FormatUint(steamLastPart*2, 10) + "]"
}

func (user *SteamUser) Steam64Hex() string {
	steamInt64, err := strconv.ParseInt(user.Steam64, 10, 64)
	if err != nil {
		logger.Errorw("[Discord] Failed to convert Steam64 to int64.", logger.Err(err))
		return ""
	}

	return fmt.Sprintf("%x", steamInt64)
}
