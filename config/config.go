package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DISCORD_API_KEY string
	GUILD_ID        string
	RCON_PASS       string
	RCON_ADDR       string
	WHITELIST_PATH  string
)

func init() {
	_ = godotenv.Load()
	DISCORD_API_KEY = must("DISCORD_API_KEY")
	GUILD_ID = must("GUILD_ID")
	RCON_PASS = must("RCON_PASS")
	RCON_ADDR = must("RCON_ADDR")
	WHITELIST_PATH = must("WHITELIST_PATH")
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
