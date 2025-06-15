/*
Filename: dynmap.go
Description: registers the commands defined by minecraft/dynmap.go
Created by: osh
        at: 17:35 on Friday, the 13th of June, 2025.
Last edited 19:04 on Sunday, the 15th of June, 2025.
*/

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

type dynmapCommand struct{}

func (dynmapCommand) Definition() *discordgo.ApplicationCommand {
	// marker worlds, marker icons and marker sets need to be fetched from the server
	// the functions to do so can result in an error, however we can just log the
	// error and proceed as normal safely, since we will always get at least an empty list back,
	// with the slash command still getting created with no available choices
	markerWorlds, err := minecraft.GetMarkerWorlds()
	if err != nil {
		log.Println("unable to retrieve dynmap worlds:", err)
	}
	// choices become autocomplete options if we have more than 25, following discord's restrictions
	markerWorldOption := func() *discordgo.ApplicationCommandOption {
		opt := &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "world",
			Description: "The dimension to add the marker to.",
			Required:    true,
		}

		if len(markerWorlds) > 25 {
			opt.Autocomplete = true
			return opt
		}

		markerWorldChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(markerWorlds))
		for i, world := range markerWorlds {
			markerWorldChoices[i] = &discordgo.ApplicationCommandOptionChoice{
				Name:  world.Title,
				Value: world.Title,
			}
		}
		opt.Choices = markerWorldChoices
		return opt
	}

	markerSets, err := minecraft.GetMarkerSets()
	if err != nil {
		log.Println("unable to retrieve marker sets:", err)
	}
	markerSetOption := func() *discordgo.ApplicationCommandOption {
		opt := &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "set",
			Description: "The marker set that this marker belongs to.",
			Required:    false,
		}

		if len(markerSets) > 25 {
			opt.Autocomplete = true
			return opt
		}

		markerSetChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(markerSets))
		for i, set := range markerSets {
			markerSetChoices[i] = &discordgo.ApplicationCommandOptionChoice{
				Name:  set.Label,
				Value: set.Name,
			}
		}
		opt.Choices = markerSetChoices
		return opt
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
					markerWorldOption(),
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "id",
						Description: "The unique id of the label, defaults to the label name.",
						Required:    false,
					},
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "icon",
						Description:  "The icon displayed for this marker on the map.",
						Required:     false,
						Autocomplete: true,
					},
					markerSetOption(),
				},
			},
		},
	}

}

func (dynmapCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	sub := i.ApplicationCommandData().Options[0]
	var msg string
	var err error
	switch sub.Name {
	case "addmarker":
		msg, err = runAddMarker(i)
	default:
		err = fmt.Errorf("unknown sub-command")
	}

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Sorry, something went wrong trying to add a marker!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func (dynmapCommand) HandleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	sub := i.ApplicationCommandData().Options[0]

	if sub.Name != "addmarker" {
		return nil
	}

	var focused *discordgo.ApplicationCommandInteractionDataOption
	for _, option := range sub.Options {
		if option.Focused {
			focused = option
			break
		}
	}

	switch focused.Name {
	case "world":
		return handleWorldAutocomplete(s, i, focused.StringValue())
	case "icon":
		return handleIconAutocomplete(s, i, focused.StringValue())
	case "set":
		return handleSetAutocomplete(s, i, focused.StringValue())
	}

	return nil
}

func handleWorldAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	markerWorlds, err := minecraft.GetMarkerWorlds()
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
	}

	content = strings.ToLower(content)
	var choices []*discordgo.ApplicationCommandOptionChoice

	for _, world := range markerWorlds {
		if strings.Contains(strings.ToLower(world.Name), content) || strings.Contains(strings.ToLower(world.Title), content) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  world.Title,
				Value: world.Title,
			})
		}

		if len(choices) >= 25 {
			break
		}
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// automatically provides autocomplete options for the users typed text using icons fetched from the server
func handleIconAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	markerIcons, err := minecraft.GetMarkerIcons()
	if err != nil {
		// no results from the server means we give back an empty set of available choices
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
	}

	content = strings.ToLower(content)
	var choices []*discordgo.ApplicationCommandOptionChoice

	for _, icon := range markerIcons {
		// since strings.Contains matches on empty string, this will also populate the choices if the user has not typed anything
		if strings.Contains(strings.ToLower(icon.Name), content) || strings.Contains(strings.ToLower(icon.Label), content) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  icon.Name,
				Value: icon.Label,
			})
		}

		// discord has a hard limit of 25 choices, will refuse to register handler if provided with more
		if len(choices) >= 25 {
			break
		}
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func handleSetAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	markerSets, err := minecraft.GetMarkerSets()
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
	}

	content = strings.ToLower(content)
	var choices []*discordgo.ApplicationCommandOptionChoice

	for _, set := range markerSets {
		if strings.Contains(strings.ToLower(set.Label), content) || strings.Contains(strings.ToLower(set.Name), content) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  set.Label,
				Value: set.Name,
			})
		}

		if len(choices) >= 25 {
			break
		}
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func runAddMarker(i *discordgo.InteractionCreate) (msg string, err error) {
	_, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	cmd := i.ApplicationCommandData()
	res, err := minecraft.AddMarker(cmd.Options)
	if err != nil {
		return
	}

	// default values if optional arguments are nonexistent for the response message
	if res.Set == "" {
		res.Set = "Markers"
	}
	if res.Icon == "" {
		res.Icon = "default"
	}

	msg = fmt.Sprintf("✅ Created marker `%v/%v` at [%v, %v, %v] on map `%v` with a %v icon.", res.Set, res.Label, res.XCoord, res.YCoord, res.ZCoord, res.WorldName, res.Icon)
	return
}

func init() {
	registry.Register(dynmapCommand{})
}
