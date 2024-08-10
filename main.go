package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/commands"
	"github.com/paranoia8972/dg-bot/config"
	events "github.com/paranoia8972/dg-bot/events/server"
	"github.com/paranoia8972/dg-bot/functions/handlers"
	"github.com/paranoia8972/dg-bot/functions/utils"
)

func main() {
	log.SetFlags(0)
	cfg := config.LoadConfig()
	utils.ConnectDB(cfg.MongoURI, cfg.DatabaseName)

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		color.Red("Error creating Discord session: %v", err)
		return
	}
	dg.AddHandler(events.MemberJoin)
	dg.AddHandler(events.InfiniteStorySend)
	dg.AddHandler(events.InfiniteStoryMessageDelete)
	dg.AddHandler(handlers.InteractionHandler)

	dg.Identify.Intents = discordgo.IntentsAll
	err = dg.Open()
	if err != nil {
		color.Red("Error opening Discord session: %v", err)
		return
	}

	commands.RegisterCommands(dg, cfg)

	err = utils.SetStatus(dg)
	if err != nil {
		log.Fatalf("Error setting status: %v", err)
	}

	color.Green("Bot is now running. Press CTRL+C to exit.")
	defer dg.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	color.Yellow("\nShutting down gracefully...")
}
