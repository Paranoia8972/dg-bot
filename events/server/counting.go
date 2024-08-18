package events

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/functions/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func CountingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	guildID := m.GuildID
	channelID := m.ChannelID
	messageContent := m.Content

	game, err := utils.GetCountingGame(guildID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			color.Red("Error retrieving counting game: %v", err)
		}
		return
	}

	if game.ChannelID != channelID {
		return
	}

	number, err := strconv.Atoi(strings.TrimSpace(messageContent))
	if err != nil {
		return
	}

	if number == game.LastNumber+1 {
		game.LastNumber = number
		err = utils.SetCountingGame(guildID, channelID, game.LastNumber)
		if err != nil {
			color.Red("Error updating counting game: %v", err)
		}
		s.MessageReactionAdd(channelID, m.ID, "✅")
	} else {
		game.LastNumber = 0
		err = utils.SetCountingGame(guildID, channelID, game.LastNumber)
		if err != nil {
			color.Red("Error resetting counting game: %v", err)
		}
		s.MessageReactionAdd(channelID, m.ID, "❌")
	}
}
