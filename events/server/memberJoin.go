package events

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/paranoia8972/dg-bot/functions/utils"
)

// TODO: Get welcome image working
func MemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	// Send welcome message
	sendWelcomeMessage(s, m.User, m.GuildID)

	rolesToAssign := []string{
		"1082708287310544937", // Member
	}

	for _, roleID := range rolesToAssign {
		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
		if err != nil {
			fmt.Printf("Error assigning role %s to user %s: %v\n", roleID, m.User.ID, err)
		}
	}
}

func sendWelcomeMessage(s *discordgo.Session, user *discordgo.User, guildID string) {
	welcomeMessage, err := utils.GetWelcomeMessage(guildID)
	if err != nil {
		fmt.Println("Error fetching welcome message:", err)
		return
	}

	welcomeMessage = replaceUserPlaceholder(welcomeMessage, user.ID)

	channelID, err := utils.GetWelcomeChannel(guildID)
	if err != nil {
		fmt.Println("Error fetching welcome channel:", err)
		return
	}

	_, err = s.ChannelMessageSend(channelID, welcomeMessage)
	if err != nil {
		fmt.Println("Error sending welcome message:", err)
	}
}

func replaceUserPlaceholder(message, ID string) string {
	return strings.ReplaceAll(message, "{user}", "<@"+ID+">")
}
