package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/smallblue2/trogbot/registry"
)

type pingCommand struct{}

func (pingCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Returns back 'Pong!' - a liveness test.",
	}
}

func (pingCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Pong!"},
	})
}

func init() {
	registry.Register(pingCommand{})
}
