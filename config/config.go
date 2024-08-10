package config

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Token        string
	MongoURI     string
	DatabaseName string
	GuildID      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	databaseName := extractDatabaseName(mongoURI)

	return &Config{
		Token:        os.Getenv("DISCORD_BOT_TOKEN"),
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		GuildID:      os.Getenv("GUILD_ID"),
	}
}

func extractDatabaseName(mongoURI string) string {
	uri, err := url.Parse(mongoURI)
	if err != nil {
		log.Fatalf("Invalid MongoDB URI: %v", err)
	}
	path := strings.Trim(uri.Path, "/")
	return path
}
