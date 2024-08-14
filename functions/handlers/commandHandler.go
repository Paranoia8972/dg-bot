package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/paranoia8972/dg-bot/commands"
)

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		commands.PingCommand(s, i)
	case "welcome":
		switch i.ApplicationCommandData().Options[0].Name {
		case "set":
			commands.WelcomeSetCommand(s, i)
		case "remove":
			commands.WelcomeRemoveCommand(s, i)
		}
	case "infinite-story":
		switch i.ApplicationCommandData().Options[0].Name {
		case "setup":
			commands.InfiniteStorySetupCommand(s, i)
		case "disable":
			commands.InfiniteStoryDisableCommand(s, i)
		}
	case "clear":
		commands.RemoveMessagesCommand(s, i)
	case "status":
		switch i.ApplicationCommandData().Options[0].Name {
		case "add":
			commands.AddStatusCommand(s, i)
		case "get":
			commands.GetStatusCommand(s, i)
		case "remove":
			commands.RemoveStatusCommand(s, i)
		}
	case "counting":
		switch i.ApplicationCommandData().Options[0].Name {
		case "setup":
			commands.CountingSetupCommand(s, i)
		case "disable":
			commands.CountingDisableCommand(s, i)

		}
	}
}
