package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/smallblue2/trogbot/commands"
	"github.com/smallblue2/trogbot/config"
	"github.com/smallblue2/trogbot/registry"
	"github.com/smallblue2/trogbot/server"

	"github.com/bwmarrin/discordgo"
)

func main() {
	dg, err := discordgo.New("Bot " + config.DISCORD_API_KEY)
	if err != nil {
		log.Fatalf("Failed to create a bot discord session: %s", err)
	}
	config.DISCORD_SESSION = dg

	// Handlers

	// On receiving the READY handler - write our slash commands
	dg.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, g := range r.Guilds {
			log.Printf("Connected to %s (%s)\n", g.Name, g.ID)
		}

		appID := s.State.User.ID
		// Register commands specifically to our guild (quicker)
		if _, err := s.ApplicationCommandBulkOverwrite(appID, config.GUILD_ID, registry.AllDefinitions()); err != nil {
			log.Fatalf("Failed to sync slash commands: %s\n", err)
		}
		// Wipe global commands (cleanup dev artefacts)
		if _, err := s.ApplicationCommandBulkOverwrite(appID, "", []*discordgo.ApplicationCommand{}); err != nil {
			log.Printf("Failed to wipe global commands: %s\n", err)
		}
		log.Println("Slash commands synced!")
	})

	// Interaction Handler for commands
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		cmdName := i.ApplicationCommandData().Name
		if cmd, ok := registry.Lookup(cmdName); ok {
			log.Printf("User '%s' ran command '%s'\n", i.Member.DisplayName(), cmdName)
			err = cmd.Run(s, i)
			if err != nil {
				log.Printf("Failed to run command '%s': %s\n", cmdName, err)
			}
		}
	})

	// Open Websocket
	if err := dg.Open(); err != nil {
		log.Fatalf("Failed to open websocket to Discord: %s", err)
	}
	log.Printf("Bot is up.")

	// Start notif server
	go func() {
		server.StartServer()
		log.Printf("Starting backup notification server on %s\n", config.NOTIF_SERVER_PORT)
		if err := http.ListenAndServe(config.NOTIF_SERVER_PORT, nil); err != nil {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	// Keep the websocket open unless SIGINT or SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	signalReceived := <-stop
	log.Printf("Received signal \"%s\", stopping.", signalReceived)
	_ = dg.Close()
}
