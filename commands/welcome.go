package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/functions/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func WelcomeSetCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	var channelID, welcomeMessage string

	for _, option := range options {
		if option.Type == discordgo.ApplicationCommandOptionSubCommand {
			for _, nestedOption := range option.Options {
				if nestedOption.Name == "channel" && nestedOption.Type == discordgo.ApplicationCommandOptionChannel {
					channel := nestedOption.ChannelValue(s)
					if channel != nil {
						channelID = channel.ID
					}
				} else if nestedOption.Name == "message" && nestedOption.Type == discordgo.ApplicationCommandOptionString {
					welcomeMessage = nestedOption.StringValue()
				}
			}
		}
	}

	guildID := i.GuildID

	if channelID != "" {
		existingChannelID, err := utils.GetWelcomeChannel(guildID)
		if err != nil && err != mongo.ErrNoDocuments {
			color.Red("Error retrieving welcome channel: %v", err)
			return
		}

		if existingChannelID != "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "A welcome channel is already set for this guild.",
				},
			})
			return
		}

		err = utils.SetWelcomeChannel(guildID, channelID)
		if err != nil {
			color.Red("Error setting welcome channel: %v", err)
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Welcome channel set successfully!",
			},
		})
	}

	if welcomeMessage != "" {
		err := utils.SetWelcomeMessage(guildID, welcomeMessage)
		if err != nil {
			color.Red("Error setting welcome message: %v", err)
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Welcome message set successfully!",
			},
		})
	}
}

func WelcomeRemoveCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var target string
	for _, option := range options {
		if option.Name == "target" && option.Type == discordgo.ApplicationCommandOptionString {
			target = option.StringValue()
		}
	}
	log.Println("target:" + target)

	guildID := i.GuildID
	var err error

	switch target {
	case "channel":
		err = utils.RemoveWelcomeChannel(guildID)
		if err != nil {
			color.Red("Error removing welcome channel: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error removing welcome channel.",
				},
			})
			return
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Welcome channel removed successfully!",
			},
		})
	case "message":
		err = utils.RemoveWelcomeMessage(guildID)
		if err != nil {
			color.Red("Error removing welcome message: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error removing welcome message.",
				},
			})
			return
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Welcome message removed successfully!",
			},
		})
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid target specified. Please choose 'message' or 'channel'.",
			},
		})
	}
}
