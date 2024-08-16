package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/functions/handlers/errorhandler"
	"github.com/paranoia8972/dg-bot/functions/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func InfiniteStorySetupCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var messageChannelID, summaryChannelID string

	for _, option := range options {
		if option.Type == discordgo.ApplicationCommandOptionSubCommand {
			for _, nestedOption := range option.Options {
				if nestedOption.Name == "message_channel_id" && nestedOption.Type == discordgo.ApplicationCommandOptionChannel {
					channel := nestedOption.ChannelValue(s)
					if channel != nil {
						messageChannelID = channel.ID
					}
				} else if nestedOption.Name == "summary_channel_id" && nestedOption.Type == discordgo.ApplicationCommandOptionChannel {
					channel := nestedOption.ChannelValue(s)
					if channel != nil {
						summaryChannelID = channel.ID
					}
				}
			}

			guildID := i.GuildID

			_, err := utils.GetInfiniteStoryChannel(guildID)
			if err == nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Infinite Story is already set up. Please disable it first before setting it up again.",
					},
				})
				return
			} else if err != mongo.ErrNoDocuments {
				color.Red("Error checking infinite story setup: %v", err)
				errorhandler.ErrorHandler(s, i, "</infinite-story setup:1272476083970904064>", err)
				return
			}

			if messageChannelID != "" && summaryChannelID != "" {
				msg, err := s.ChannelMessageSend(summaryChannelID, "# Infinite Story Summary")
				if err != nil {
					color.Red("Error sending summary message: %v", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Failed to send summary message.",
						},
					})
					return
				}

				err = utils.SetInfiniteStoryChannel(guildID, messageChannelID, summaryChannelID, "# Infinite Story Summary\n", msg.ID)
				if err != nil {
					color.Red("Error saving infinite story channel details: %v", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Failed to save infinite story channel details.",
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Infinite Story channels have been successfully set!",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to set Infinite Story channels. Please ensure both channels are specified.",
					},
				})
			}
		}
	}
}

func InfiniteStoryDisableCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID
	err := utils.RemoveInfiniteStoryChannel(guildID)
	if err != nil {
		color.Red("Error removing infinite story channel: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to remove infinite story channel.",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Infinite Story channel has been successfully removed!",
		},
	})
}
