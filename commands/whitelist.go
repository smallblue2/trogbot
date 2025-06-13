package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/smallblue2/trogbot/minecraft"
	"github.com/smallblue2/trogbot/registry"
)

type whitelistCommand struct{}

func (whitelistCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "whitelist",
		Description: "Manage the Minecraft server's whitelist",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "Show currently whitelisted players",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add a new player to the whitelist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "player",
						Description: "Minecraft username",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a player from the whitelist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "player",
						Description: "Minecraft username",
						Required:    true,
					},
				},
			},
		},
	}
}

func (whitelistCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	sub := i.ApplicationCommandData().Options[0]
	switch sub.Name {
	case "list":
		return runList(s, i)
	case "add":
		return runAdd(s, i, sub.GetOption("player").StringValue())
	case "remove":
		return runRemove(s, i, sub.GetOption("player").StringValue())
	default:
		return fmt.Errorf("unknown sub-command")
	}
}

func runList(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	list, err := minecraft.Load()
	if err != nil {
		return err
	}
	names := make([]string, len(list))
	for idx, e := range list {
		names[idx] = " - " + e.Name + "\n"
	}
	content := "Whitelisted players:\n" + strings.Join(names, "")

	log.Printf("User '%s' listed the whitelist.\n", i.Member.DisplayName())

	return respond(s, i, content)
}

func runAdd(s *discordgo.Session, i *discordgo.InteractionCreate, player string) error {

	// Validate against Mojang
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	uuid, err := minecraft.FetchUUID(ctx, player)
	if err != nil {
		return respond(s, i, fmt.Sprintf("❌ %s", err))
	}

	// Update the list
	list, _ := minecraft.Load()
	for _, e := range list {
		if strings.EqualFold(e.Name, player) {
			return respond(s, i, "✅ already whitelisted")
		}
	}

	list = append(list, minecraft.Entry{UUID: uuid, Name: player})
	if err := minecraft.Save(list); err != nil {
		return err
	}

	log.Printf("User '%s' added '%s' (%s) to the whitelist.\n", i.Member.DisplayName(), player, uuid)

	if _, err := minecraft.Exec("whitelist reload"); err != nil {
		log.Printf("Failed RCON call: %s\n", err)
		return respond(s, i, "⚠️ Saved, but failed to reload the whitelist")
	}

	return respond(s, i, "✅ added "+player)
}

func runRemove(s *discordgo.Session, i *discordgo.InteractionCreate, player string) error {
	list, _ := minecraft.Load()
	out := list[:0]
	removed := false
	for _, e := range list {
		if !strings.EqualFold(e.Name, player) {
			out = append(out, e)
		} else {
			removed = true
		}
	}
	if !removed {
		return respond(s, i, "❌ "+player+" is not on whitelist")
	}
	if err := minecraft.Save(out); err != nil {
		return err
	}

	log.Printf("User '%s' removed '%s' from the whitelist.\n", i.Member.DisplayName(), player)

	if _, err := minecraft.Exec("whitelist reload"); err != nil {
		log.Printf("Failed RCON call: %s\n", err)
		return respond(s, i, "⚠️ Saved, but failed to reload the whitelist")
	}

	return respond(s, i, "✅ Removed "+player)
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}

func init() {
	registry.Register(whitelistCommand{})
}
