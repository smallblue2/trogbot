/*
Filename: dynmap.go
Description: interfaces with dynmap marker commands to allow for the
             creation, modification and deletion of dynmap markers
Created by: osh
        at: 15:39 on Friday, the 13th of June, 2025.
Last edited 13:27 on Sunday, the 15th of June, 2025.
*/

package minecraft

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Marker struct {
	Label     string
	Id        string
	Icon      string
	Set       string
	XCoord    int
	YCoord    int
	ZCoord    int
	WorldName string
}

type MarkerIcon struct {
	Name    string
	Label   string // unsure if Name and Label are different, it might just be the same field listed twice
	Builtin bool
}

type MarkerSet struct {
	Name        string
	Label       string
	Hide        bool
	Priority    int
	DefaultIcon string
	Persistent  bool
}

type MarkerWorld struct {
	Name         string
	Loaded       bool
	Enabled      bool
	Title        string
	Center       string
	ExtraZoomOut int
	SendHealth   bool
	SendPosition bool
	Protected    bool
	ShowBorder   bool
}

// builds the `/dmarker add` command, allowing for optional fields
func execDmarkerAdd(marker Marker) (result string, err error) {
	// we're not guaranteed that the order from the slash command is correct, so just build each sequentially
	var parts []string

	parts = append(parts, "dmarker add")

	// optional id comes first
	if marker.Id != "" {
		parts = append(parts, fmt.Sprintf("id:%v", marker.Id))
	}

	// label is required and can contain spaces (that must be quoted)
	if strings.Contains(marker.Label, " ") {
		parts = append(parts, fmt.Sprintf("\"%v\"", marker.Label))
	} else {
		parts = append(parts, marker.Label)
	}

	if marker.Icon != "" {
		parts = append(parts, fmt.Sprintf("icon:%v", marker.Icon))
	}

	if marker.Set != "" {
		parts = append(parts, fmt.Sprintf("set:%v", marker.Set))
	}

	// required world coords and name
	parts = append(parts, fmt.Sprintf("x:%v", marker.XCoord))
	parts = append(parts, fmt.Sprintf("y:%v", marker.YCoord))
	parts = append(parts, fmt.Sprintf("z:%v", marker.ZCoord))
	parts = append(parts, fmt.Sprintf("world:%v", marker.WorldName))

	cmd := strings.Join(parts, " ")
	return Exec(cmd)
}

// adds a marker given marker fields
func AddMarker(slashCommandData []*discordgo.ApplicationCommandInteractionDataOption) (marker Marker, err error) {
	if len(slashCommandData) != 1 {
		err = fmt.Errorf("malformed application command")
		return
	}

	// option content can be given in an arbitrary order, with optional arguments arbitrarily filled
	// save as a Marker by looking at the given name of the slash command field
	res := slashCommandData[0]
	resOptions := res.Options
	for _, option := range resOptions {
		switch option.Name {
		case "label":
			marker.Label = option.StringValue()
		case "id":
			marker.Id = option.StringValue()
		case "icon":
			marker.Icon = option.StringValue()
		case "set":
			marker.Set = option.StringValue()
		case "x":
			marker.XCoord = int(option.IntValue())
		case "y":
			marker.YCoord = int(option.IntValue())
		case "z":
			marker.ZCoord = int(option.IntValue())
		case "world":
			marker.WorldName = option.StringValue()
		default:
			err = fmt.Errorf("unrecognised option %v", option.Name)
			return
		}
	}

	return
}

// retrieves the currently available worlds using `/dmap worldlist` to populat ethe slash command options
// looks like (newline seperated):
//
// world <world id>: loaded=<bool>, enabled=<bool>, title=<the dimension name>, center=<x.x/y.y/z.z>, extrazoomout=<int>, sendpositon=<bool>, protected=<bool>, showborder=<bool>
// eg:
// world world_172429123: loaded=true, enabled=true, title=world, center=-32.0/70.0/-80.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=true
func parseMarkerWorlds(result string) (markerWorlds []MarkerWorld, err error) {
	re := regexp.MustCompile(`world ([\w\-]+): loaded=(true|false), enabled=(true|false), title=(\w+), center=([\d\.\-/]+), extrazoomout=(\d+), sendhealth=(true|false), sendposition=(true|false), protected=(true|false), showborder=(true|false)`)
	matches := re.FindAllStringSubmatch(result, -1)
	var extraZoomOut int
	for _, match := range matches {
		extraZoomOut, err = strconv.Atoi(match[6])
		if err != nil {
			return
		}

		markerWorlds = append(markerWorlds, MarkerWorld{
			Name:         match[1],
			Loaded:       match[2] == "true",
			Enabled:      match[3] == "true",
			Title:        match[4],
			Center:       match[5],
			ExtraZoomOut: extraZoomOut,
			SendHealth:   match[7] == "true",
			SendPosition: match[8] == "true",
			Protected:    match[9] == "true",
			ShowBorder:   match[10] == "true",
		})
	}

	return
}

func GetMarkerWorlds() (markerWorlds []MarkerWorld, err error) {
	result, err := Exec("dmap worldlist")
	if err != nil {
		return
	}

	return parseMarkerWorlds(result)
}

// retrieves currently available marker icons using `/dmarker icons` to populate the slash command options
// looks like (newline seperated):
//
// <label>: label:"<label>", builtin:<bool>
// eg:
// anchor: label:"anchor", builtin:true
func parseMarkerIcons(result string) (markerIcons []MarkerIcon, err error) {
	re := regexp.MustCompile(`(\w+): label:"([^"]+)", builtin:(true|false)`)
	matches := re.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		markerIcons = append(markerIcons, MarkerIcon{
			Name:    match[1],
			Label:   match[2],
			Builtin: match[3] == "true",
		})
	}

	return
}

func GetMarkerIcons() (markerIcons []MarkerIcon, err error) {
	result, err := Exec("dmarker icons")
	if err != nil {
		return
	}

	return parseMarkerIcons(result)
}

// retrieves currently available marker sets using `/dmarker listsets` to populate the slash command options
// looks like (newline seperated):
//
// <marker set name>: label:"<label display name>", hide:<bool>, prio:<int>, deficon:<default icon name>
// eg:
// markers: label:"Markers", hide:false, prio:0, deficon:default
func parseMarkerSets(result string) (markerSets []MarkerSet, err error) {
	re := regexp.MustCompile(`(\w+): label:"([^"]+)", hide:(true|false), prio:(\d+), deficon:(\w+)(?:, persistent=(true|false))?`)
	matches := re.FindAllStringSubmatch(result, -1)
	var priority int
	for _, match := range matches {
		priority, err = strconv.Atoi(match[4])
		if err != nil {
			return
		}

		// persistent field is optional
		var persistent bool
		if len(match) > 6 && match[6] != "" {
			persistent = match[6] == "true"
		}

		markerSets = append(markerSets, MarkerSet{
			Name:        match[1],
			Label:       match[2],
			Hide:        match[3] == "true",
			Priority:    priority,
			DefaultIcon: match[5],
			Persistent:  persistent,
		})
	}

	return
}

func GetMarkerSets() (markerSets []MarkerSet, err error) {
	result, err := Exec("dmarker listsets")
	if err != nil {
		return
	}

	return parseMarkerSets(result)
}
