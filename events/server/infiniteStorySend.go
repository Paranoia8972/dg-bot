package events

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/paranoia8972/dg-bot/functions/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func InfiniteStorySend(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID
	channelID := m.ChannelID
	messageContent := m.Content
	if len(strings.Fields(messageContent)) > 1 && !m.Author.Bot {
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			color.Red("Error deleting message: %v", err)
			return
		}
	}

	infiniteStoryChannel, err := utils.GetInfiniteStoryChannel(guildID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			color.Red("Error retrieving infinite story channel: %v", err)
		}
		return
	}

	if infiniteStoryChannel.MessageChannelID == channelID {
		if len(strings.Fields(messageContent)) == 1 {
			infiniteStoryChannel.Summary += " " + strings.TrimSpace(messageContent)
		}

		if infiniteStoryChannel.MessageID == "" {
			msg, err := s.ChannelMessageSend(infiniteStoryChannel.SummaryChannelID, infiniteStoryChannel.Summary)
			if err != nil {
				color.Red("Error sending new summary message: %v", err)
				return
			}
			infiniteStoryChannel.MessageID = msg.ID
			err = utils.SetInfiniteStoryChannel(guildID, infiniteStoryChannel.MessageChannelID, infiniteStoryChannel.SummaryChannelID, infiniteStoryChannel.Summary, infiniteStoryChannel.MessageID)
			if err != nil {
				color.Red("Error saving updated summary: %v", err)
			}
			return
		}
	}
}

func InfiniteStoryMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	guildID := m.GuildID
	messageID := m.ID

	infiniteStoryChannel, err := utils.GetInfiniteStoryChannel(guildID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			color.Red("Error retrieving infinite story channel: %v", err)
		}
		return
	}

	if infiniteStoryChannel.MessageID == messageID {
		msg, err := s.ChannelMessageSend(infiniteStoryChannel.SummaryChannelID, infiniteStoryChannel.Summary)
		if err != nil {
			color.Red("Error sending new summary message: %v", err)
			return
		}
		infiniteStoryChannel.MessageID = msg.ID
		err = utils.SetInfiniteStoryChannel(guildID, infiniteStoryChannel.MessageChannelID, infiniteStoryChannel.SummaryChannelID, infiniteStoryChannel.Summary, infiniteStoryChannel.MessageID)
		if err != nil {
			color.Red("Error saving updated summary: %v", err)
		}
	}
}
