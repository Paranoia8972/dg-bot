package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InfiniteStoryChannel struct {
	GuildID          string `bson:"guild_id"`
	MessageChannelID string `bson:"message_channel_id"`
	SummaryChannelID string `bson:"summary_channel_id"`
	Summary          string `bson:"summary"`
	MessageID        string `bson:"message_id"`
}

func SetInfiniteStoryChannel(guildID, messageChannelID, summaryChannelID, summary, messageID string) error {
	collection := Client.Database("discordgo").Collection("infinite_story_channels")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"guild_id": guildID},
		bson.M{"$set": bson.M{
			"message_channel_id": messageChannelID,
			"summary_channel_id": summaryChannelID,
			"summary":            summary,
			"message_id":         messageID,
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

func GetInfiniteStoryChannel(guildID string) (InfiniteStoryChannel, error) {
	collection := Client.Database("discordgo").Collection("infinite_story_channels")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result InfiniteStoryChannel
	err := collection.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&result)
	return result, err
}

func RemoveInfiniteStoryChannel(guildID string) error {
	collection := Client.Database("discordgo").Collection("infinite_story_channels")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"guild_id": guildID})
	return err
}
