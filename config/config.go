package config

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	DISCORD_API_KEY   string
	GUILD_ID          string
	BOT_CHANNEL_ID    string
	RCON_PASS         string
	RCON_ADDR         string
	WHITELIST_PATH    string
	DISCORD_SESSION   *discordgo.Session
	NOTIF_SERVER_PORT string
)

func init() {
	_ = godotenv.Load()
	DISCORD_API_KEY = must("DISCORD_API_KEY")
	GUILD_ID = must("GUILD_ID")
	BOT_CHANNEL_ID = must("BOT_CHANNEL_ID")
	RCON_PASS = must("RCON_PASS")
	RCON_ADDR = must("RCON_ADDR")
	WHITELIST_PATH = must("WHITELIST_PATH")
	NOTIF_SERVER_PORT = must("NOTIF_SERVER_PORT")
}

func must(key string) string {
	v := os.Getenv(key)
	if v == "" {
		// Warn, but allow continued running for tests
		log.Printf("Env '%s' is required and not set!\n", key)
		os.Setenv(key, "UNSET")
	}
	return v
}
