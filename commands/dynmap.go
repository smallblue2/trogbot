/*
Filename: dynmap.go
Description: registers the commands defined by minecraft/dynmap.go
Created by: osh
        at: 17:35 on Friday, the 13th of June, 2025.
Last edited 19:48 on Friday, the 13th of June, 2025.
*/

package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/smallblue2/trogbot/minecraft"
	"github.com/smallblue2/trogbot/registry"
)

type dynmapCommand struct{}

func (dynmapCommand) Definition() *discordgo.ApplicationCommand {
	// marker worlds, marker icons and marker sets need to be fetched from the server,
	// but we have no error propagation so we log that no choices were set on error
	// this is safe to do so because we will always get at least an empty list back

	markerWorlds, err := minecraft.GetMarkerWorlds()
	if err != nil {
		log.Println("unable to retrieve dynmap worlds:", err)
	}
	markerWorldChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(markerWorlds))
	for i, world := range markerWorlds {
		markerWorldChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  world.Title,
			Value: world.Title,
		}
	}

	markerIcons, err := minecraft.GetMarkerIcons()
	if err != nil {
		log.Println("unable to retrieve marker icons:", err)
	}
	markerIconChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(markerIcons))
	for i, icon := range markerIcons {
		markerIconChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  icon.LabelName,
			Value: icon.Label,
		}
	}

	markerSets, err := minecraft.GetMarkerSets()
	if err != nil {
		log.Println("unable to retrieve marker sets:", err)
	}
	markerSetChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(markerSets))
	for i, set := range markerSets {
		markerSetChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  set.Label,
			Value: set.Name,
		}
	}

	return &discordgo.ApplicationCommand{
		Name:        "dynmap",
		Description: "Dynmap commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "addmarker",
				Description: "Adds a marker to the map.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "label",
						Description: "The label that gets displayed on the map.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "x",
						Description: "The x-coordinate on the map to create the marker at.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "y",
						Description: "The y-coordinate on the map to create the marker at.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "z",
						Description: "The z-coordinate on the map to create the marker at.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "world",
						Description: "The dimension to add the marker to.",
						Required:    true,
						Choices:     markerWorldChoices,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "id",
						Description: "The unique id of the label, defaults to the label name.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "icon",
						Description: "The icon displayed for this marker on the map.",
						Required:    false,
						Choices:     markerIconChoices,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "set",
						Description: "The marker set that this marker belongs to.",
						Required:    false,
						Choices:     markerSetChoices,
					},
				},
			},
		},
	}

}

func (dynmapCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	sub := i.ApplicationCommandData().Options[0]
	switch sub.Name {
	case "addmarker":
		return runAddMarker(s, i)
	default:
		return fmt.Errorf("unknown sub-command")
	}
}

func runAddMarker(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	return nil
}

func init() {
	registry.Register(dynmapCommand{})
}
