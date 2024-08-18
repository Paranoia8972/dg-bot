package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/functions/utils"
)

func CountingSetupCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID
	channelID := i.ChannelID

	err := utils.SetCountingGame(guildID, channelID, 0)
	if err != nil {
		color.Red("Error setting up counting game: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set up counting game.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Counting game has been successfully set up!",
		},
	})
}

func CountingDisableCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID

	err := utils.RemoveCountingGame(guildID)
	if err != nil {
		color.Red("Error disabling counting game: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to disable counting game.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Counting game has been successfully disabled!",
		},
	})
}
