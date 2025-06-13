package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/smallblue2/trogbot/minecraft"
	"github.com/smallblue2/trogbot/registry"
)

type listPlayersCommands struct{}

func (listPlayersCommands) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "listplayers",
		Description: "List the players that are currently online",
	}
}

func (listPlayersCommands) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	result, err := minecraft.Exec("list")
	if err != nil {
		return err
	}

	msg := minecraft.GetOnlinePlayerMsg(result)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}

func init() {
	registry.Register(listPlayersCommands{})
}
