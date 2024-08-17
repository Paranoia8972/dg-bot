package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/paranoia8972/dg-bot/functions/utils"
	"github.com/paranoia8972/dg-bot/schemas"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	status := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	client := utils.Client
	if client == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "MongoDB client is not initialized",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("discordgo").Collection("botstatus")

	// Check if the status already exists
	var existingStatus schemas.BotStatus
	err := collection.FindOne(ctx, bson.M{"status": status}).Decode(&existingStatus)
	if err == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Status already exists!",
			},
		})
		return
	}

	_, err = collection.InsertOne(ctx, schemas.BotStatus{Status: status})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to add status: %v", err),
			},
		})
		return
	}

	if err := utils.SetStatus(s); err != nil {
		log.Printf("Failed to update status: %v", err)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Status added successfully!",
		},
	})
}

func GetStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	client := utils.Client
	if client == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "MongoDB client is not initialized",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("discordgo").Collection("botstatus")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to retrieve statuses: %v", err),
			},
		})
		return
	}
	defer cursor.Close(ctx)

	var statuses []schemas.BotStatus
	if err = cursor.All(ctx, &statuses); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to decode statuses: %v", err),
			},
		})
		return
	}

	response := "Statuses:\n"
	for _, status := range statuses {
		response += fmt.Sprintf("- %s (ID: %s)\n", status.Status, status.ID.Hex())
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   1 << 6,
		},
	})
}

func RemoveStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	statusID := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	client := utils.Client
	if client == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "MongoDB client is not initialized",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("discordgo").Collection("botstatus")
	id, err := primitive.ObjectIDFromHex(statusID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Invalid status ID: %v", err),
			},
		})
		return
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to remove status: %v", err),
			},
		})
		return
	}

	if err := utils.SetStatus(s); err != nil {
		log.Printf("Failed to update status: %v", err)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Status removed successfully!",
		},
	})
}
