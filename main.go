package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/matthewpi/snaily/backend"
	"github.com/matthewpi/snaily/bot"
	"github.com/matthewpi/snaily/command"
	"github.com/matthewpi/snaily/command/commands"
	"github.com/matthewpi/snaily/config"
	"github.com/matthewpi/snaily/dca"
	"github.com/matthewpi/snaily/events"
	"github.com/matthewpi/snaily/logger"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	buildVersion string
	buildBranch  string
	buildCommit  string
	buildDate    string
)

func init() {
	if buildVersion == "" {
		buildVersion = "v0.0.1"
	}

	if buildBranch == "" {
		buildBranch = "master"
	}

	if buildCommit == "" {
		buildCommit = "unknown"
	}

	if buildDate == "" {
		buildDate = "2000-01-01"
	}
}

func main() {
	log.SetFlags(0)

	log.Println(`   ____          _ __`)
	log.Println(`  / __/__  ___ _(_) /_ __`)
	log.Println(` _\ \/ _ \/ _  / / / // /`)
	log.Println(`/___/_//_/\_,_/_/_/\_, /`)
	log.Println(`                  /___/`)
	log.Printf("  %s    %s/%s\n\n", buildVersion, buildBranch, buildCommit)

	if err := logger.Initialize(); err != nil {
		log.Fatalf("[Preflight] Failed to initialize logger: %v", err)
		return
	}

	if runtime.GOARCH != "amd64" {
		logger.Fatal("[Preflight] Access only supports 'amd64' systems.")
		return
	}

	if runtime.GOOS != "linux" {
		logger.Fatal("[Preflight] Access only supports 'Linux' operating systems.")
		return
	}

	if err := config.Load("Snaily", buildVersion, buildBranch, buildCommit, buildDate); err != nil {
		logger.Fatalw("[Preflight] Failed to load configuration.", logger.Err(err))
		return
	}
	logger.Debug("[Preflight] All preflight tasks were successful.")

	redis := &backend.RedisDriver{}
	if err := redis.Connect(config.Get().Backend.Redis.URI, config.Get().Backend.Redis.Password, config.Get().Backend.Redis.Database); err != nil {
		logger.Fatalw("[Redis] Failed to connect to remote server.", logger.Err(err))
		return
	}
	logger.Info("[Redis] Connected to remote server.")

	bot.SetBot(&bot.Bot{
		Config: config.Get(),
		Commands: []*command.Command{
			// Base Commands
			commands.Hypixel(),
			commands.Info(),
			commands.Ping(),
			commands.Steam(),

			// Moderation Commands
			commands.Ban(),
			commands.Kick(),
			commands.Mute(),
			commands.Purge(),

			// Music Commands
			commands.Pause(),
			commands.Play(),
			commands.Playing(),
			commands.Playlist(),
			commands.Queue(),
			commands.Shuffle(),
			commands.Skip(),
			commands.Stop(),
		},
		Redis:       redis,
		Queue:       map[string][]*bot.Request{},
		Playing:     map[string]*bot.Request{},
		Shuffle:     map[string]bool{},
		MusicStream: map[string]*dca.StreamingSession{},
	})

	discord, err := discordgo.New("Bot " + config.Get().Discord.Token)
	if err != nil {
		logger.Fatalw("[Discord] Failed to authenticate bot user.", logger.Err(err))
		return
	}

	discord.AddHandler(events.ReadyEvent)
	discord.AddHandler(events.MessageCreateEvent)
	discord.AddHandler(events.MessageUpdateEvent)
	discord.AddHandler(events.MessageDeleteEvent)
	discord.AddHandler(events.GuildRoleCreateEvent)
	discord.AddHandler(events.GuildRoleDeleteEvent)

	err = discord.Open()
	if err != nil {
		logger.Fatalw("[Discord] Failed to connect to discord.", logger.Err(err))
		return
	}
	bot.GetBot().Session = discord

	botUser, err := discord.User("@me")
	if err != nil {
		logger.Fatalw("[Discord] Failed to obtain account details.", logger.Err(err))
		return
	}
	bot.GetBot().User = botUser

	// Start the music thread.
	for key := range bot.GetBot().Config.Discord.Guilds {
		bot.GetBot().Shuffle[key] = false
		bot.GetBot().Music(key)
	}

	buildVersion = ""
	buildBranch = ""
	buildCommit = ""
	buildDate = ""

	logger.Infof("[Discord] Connected successfully, bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	logger.Warn("[System] Received signal from the system, closing connections.")
	_ = discord.Close()
	_ = redis.Client.Close()
}
