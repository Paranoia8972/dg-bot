package utils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

type BotStatus struct {
	Status string `bson:"status"`
}

var DefaultStatus = BotStatus{
	Status: "online",
}

func SetStatus(s *discordgo.Session) error {
	client := Client
	if client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("discordgo").Collection("botstatus")
	var statuses []BotStatus
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to find statuses: %v", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &statuses); err != nil {
		return fmt.Errorf("failed to decode statuses: %v", err)
	}

	if len(statuses) == 0 {
		return fmt.Errorf("no statuses found")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(statuses), func(i, j int) {
		statuses[i], statuses[j] = statuses[j], statuses[i]
	})

	err = s.UpdateGameStatus(0, statuses[0].Status)
	if err != nil {
		return fmt.Errorf("failed to set status: %v", err)
	}

	ticker := time.NewTicker(10 * time.Second)
	lastStatus := statuses[0].Status
	go func() {
		for range ticker.C {
			var newStatus string
			for {
				r.Shuffle(len(statuses), func(i, j int) {
					statuses[i], statuses[j] = statuses[j], statuses[i]
				})
				newStatus = statuses[0].Status
				if newStatus != lastStatus {
					break
				}
			}
			err := s.UpdateGameStatus(0, newStatus)
			if err != nil {
				log.Println("Failed to set status:", err)
			} else {
				lastStatus = newStatus
			}
		}
	}()

	return nil
}
