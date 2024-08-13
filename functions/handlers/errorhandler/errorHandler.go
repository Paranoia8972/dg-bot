package errorhandler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func ErrorHandler(s *discordgo.Session, i *discordgo.InteractionCreate, commandName string, err error) {
	color.Red("An error occurred in command %s: %v", commandName, err)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "An error occurred while processing the `" + commandName + "` command. Please contact <@982984144567017493> if this issue persists.",
		},
	})
}
