package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CountingGame struct {
	GuildID    string `bson:"guild_id"`
	ChannelID  string `bson:"channel_id"`
	LastNumber int    `bson:"last_number"`
}

var countingGameCache = make(map[string]CountingGame)

func GetCountingGame(guildID string) (CountingGame, error) {
	if game, found := countingGameCache[guildID]; found {
		return game, nil
	}
	collection := Client.Database("discordgo").Collection("counting_games")
	collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"guild_id": 1},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var game CountingGame
	err := collection.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&game)
	if err == nil {
		countingGameCache[guildID] = game
	}
	return game, err
}

func SetCountingGame(guildID, channelID string, lastNumber int) error {
	collection := Client.Database("discordgo").Collection("counting_games")
	collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"guild_id": 1},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"guild_id": guildID},
		bson.M{"$set": bson.M{
			"channel_id":  channelID,
			"last_number": lastNumber,
		}},
		options.Update().SetUpsert(true),
	)
	if err == nil {
		countingGameCache[guildID] = CountingGame{
			GuildID:    guildID,
			ChannelID:  channelID,
			LastNumber: lastNumber,
		}
	}
	return err
}

func RemoveCountingGame(guildID string) error {
	collection := Client.Database("discordgo").Collection("counting_games")
	collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"guild_id": 1},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"guild_id": guildID})
	return err
}
