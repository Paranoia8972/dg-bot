package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/config"
)

func RegisterCommands(dg *discordgo.Session, cfg *config.Config) {
	one := float64(1)
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "welcome",
			Description: "Manage welcome settings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set",
					Description: "Set the welcome channel or message",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "The channel to set as the welcome channel",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "message",
							Description: "The welcome message to set. Use {user} to mention the user",
							Required:    false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "Remove the welcome channel and message",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "target",
							Description: "Specify what to remove (message or channel)",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "message",
									Value: "message",
								},
								{
									Name:  "channel",
									Value: "channel",
								},
							},
						},
					},
				},
			},
		},
		{
			Name:        "clear",
			Description: "Remove a certain number of messages optionally from a specific user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "The number of messages to remove",
					Required:    true,
					MinValue:    &one,
					MaxValue:    99,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user whose messages to remove",
					Required:    false,
				},
			},
		},
		{
			Name:        "infinite-story",
			Description: "Manage infinite story settings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "setup",
					Description: "Set the infinite story channel and summary",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "message_channel_id",
							Description: "The channel to set as the infinite story channel",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "summary_channel_id",
							Description: "The initial summary message",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "disable",
					Description: "Disable the infinite story channel",
				},
			},
		},
		{
			Name:        "ping",
			Description: "Get the bot's latency",
		},
		{
			Name:        "status",
			Description: "Manage bot statuses",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Add a new status",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "status",
							Description: "The status to add",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "get",
					Description: "Get all statuses",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "remove",
					Description: "Remove a status",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "status",
							Description: "The status to remove",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "counting",
			Description: "Manage counting settings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "setup",
					Description: "Set the counting channel",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "The channel to set as the counting channel",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "disable",
					Description: "Disable the counting channel",
				},
			},
		},
	}

	existingCommands, err := dg.ApplicationCommands(dg.State.User.ID, cfg.GuildID)
	if err != nil {
		color.Red("Cannot fetch existing commands: %v", err)
		return
	}
	for _, cmd := range existingCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, cfg.GuildID, cmd.ID)
		if err != nil {
			color.Red("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
	for _, cmd := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, cfg.GuildID, cmd)
		if err != nil {
			color.Red("Cannot create '%v' command: %v", cmd.Name, err)
		} else {
			color.Blue("Command '%v' created successfully", cmd.Name)
		}
	}
}
