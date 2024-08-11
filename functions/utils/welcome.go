package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WelcomeChannel struct {
	GuildID   string `bson:"guild_id"`
	ChannelID string `bson:"channel_id"`
}

func SetWelcomeMessage(guildID, message string) error {
	collection := Client.Database("discordgo").Collection("welcome_messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"guild_id": guildID},
		bson.M{"$set": bson.M{"message": message}},
		options.Update().SetUpsert(true),
	)
	return err
}

func GetWelcomeMessage(guildID string) (string, error) {
	collection := Client.Database("discordgo").Collection("welcome_messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Message string `bson:"message"`
	}
	err := collection.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Message, nil
}

func RemoveWelcomeMessage(guildID string) error {
	collection := Client.Database("discordgo").Collection("welcome_messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"guild_id": guildID})
	return err
}

func SetWelcomeChannel(guildID, channelID string) error {
	collection := Client.Database("discordgo").Collection("welcome_channels")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"guild_id": guildID},
		bson.M{"$set": bson.M{"channel_id": channelID}},
		options.Update().SetUpsert(true),
	)
	return err
}

func RemoveWelcomeChannel(guildID string) error {
	collection := Client.Database("discordgo").Collection("welcome_channels")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"guild_id": guildID})
	return err
}

func GetWelcomeChannel(guildID string) (string, error) {
	collection := Client.Database("discordgo").Collection("welcome_channels")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result WelcomeChannel
	err := collection.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.ChannelID, nil
}
