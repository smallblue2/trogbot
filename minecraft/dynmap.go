/*
Filename: dynmap.go
Description: interfaces with dynmap marker commands to allow for the
             creation, modification and deletion of dynmap markers
Created by: osh
        at: 15:39 on Friday, the 13th of June, 2025.
Last edited 17:40 on Saturday, the 14th of June, 2025.
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

// // builds the `/dmarker add` command, allowing for optional fields
func execDmarkerAdd(marker Marker) (result string, err error) {
	// we're not guaranteed that the order from the slash command is correct, so just build each sequentially
	var parts []string

	parts = append(parts, "dmarker add")

	// optional id first
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
func execDmapWorldlist() (result string, err error) {
	// result = `
	// 		 world world: loaded=true, enabled=true, title=world, center=-32.0/70.0/-80.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=true
	// 		 world world_nether: loaded=true, enabled=true, title=the_nether, center=-32.0/70.0/-80.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=true
	// 		 world world_the_end: loaded=true, enabled=true, title=the_end, center=-32.0/70.0/-80.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=true
	// 		 `

	result, err = minecraft.Exec("dmap worldlist")
	return
}

func GetMarkerWorlds() (markerWorlds []MarkerWorld, err error) {
	result, err := execDmapWorldlist()
	if err != nil {
		return
	}

	re := regexp.MustCompile(`world (\w+): loaded=(\w+), enabled=(\w+), title=(\w+), center=([\d\.\-/]+), extrazoomout=(\d+), sendhealth=(\w+), sendposition=(\w+), protected=(\w+), showborder=(\w+)`)
	matches := re.FindAllStringSubmatch(result, -1)
	var extraZoomOut int
	for _, match := range matches {
		extraZoomOut, err = strconv.Atoi(match[6])
		if err != nil {
			return
		}

		// just check that the bools are actually bools
		for _, m := range []string{match[2], match[3], match[7], match[8], match[9], match[10]} {
			if m != "true" && m != "false" {
				return markerWorlds, fmt.Errorf("unable to parse %v as bool\n", match[3])
			}
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

// retrieves currently available marker icons using `/dmarker icons` to populate the slash command options
// looks like (newline seperated):
//
// <label>: label:"<label>", builtin:<bool>
// eg:
// anchor: label:"anchor", builtin:true
func execDmarkerIcons() (result string, err error) {
	// result = `
	// 		 anchor: label:"anchor", builtin:true
	// 		 basket: label:"basket", builtin:true
	// 		 bed: label:"bed", builtin:true
	// 		 brick: label:"brick", builtin:true
	// 		 compass: label:"compass", builtin:true
	// 		 default: label:"default", builtin:true
	// 		 diamond: label:"diamond", builtin:true
	// 		 house: label:"house", builtin:true
	// 		 lightbulb: label:"lightbulb", builtin:true
	// 		 offlineuser: label:"offlineuser", builtin:true
	// 		 pin: label:"pin", builtin:true
	// 		 redbasket: label:"redbasket", builtin:true
	// 		 sign: label:"sign", builtin:true
	// 		 skull: label:"skull", builtin:true
	// 		 tower: label:"tower", builtin:true
	// 		 tree: label:"tree", builtin:true
	// 		 world: label:"world", builtin:true
	// 		 custom_shop: label:"Custom Shop Icon", builtin:false
	// 		 `
	result, err = minecraft.Exec("dmarker icons")

	return
}

func GetMarkerIcons() (markerIcons []MarkerIcon, err error) {
	result, err := execDmarkerIcons()
	if err != nil {
		return
	}

	re := regexp.MustCompile(`(\w+): label:"([^"]+)", builtin:(\w+)`)
	matches := re.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		if match[3] != "true" && match[3] != "false" {
			return markerIcons, fmt.Errorf("unable to parse %v as bool\n", match[3])
		}

		markerIcons = append(markerIcons, MarkerIcon{
			Name:    match[1],
			Label:   match[2],
			Builtin: match[3] == "true",
		})
	}

	return
}

// retrieves currently available marker sets using `/dmarker listsets` to populate the slash command options
// looks like (newline seperated):
//
// <marker set name>: label:"<label display name>", hide:<bool>, prio:<int>, deficon:<default icon name>
// eg:
// markers: label:"Markers", hide:false, prio:0, deficon:default
func execDmarkerListsets() (result string, err error) {
	// result = `
	// 		 markers: label:"Markers", hide:false, prio:0, deficon:default
	// 		 shops: label:"Shop Locations", hide:false, prio:10, deficon:basket
	// 		 bases: label:"Player Bases", hide:false, prio:5, deficon:house
	// 		 farms: label:"Farms", hide:true, prio:3, deficon:tree
	// 		 warps: label:"Warp Points", hide:false, prio:15, deficon:compass
	// 		 `
	result, err = minecraft.Exec("dmarker listsets")

	return
}

func GetMarkerSets() (markerSets []MarkerSet, err error) {
	result, err := execDmarkerListsets()
	if err != nil {
		return
	}

	re := regexp.MustCompile(`(\w+): label:"([^"]+)", hide:(\w+), prio:(\d+), deficon:(\w+)`)
	matches := re.FindAllStringSubmatch(result, -1)
	var priority int
	for _, match := range matches {
		priority, err = strconv.Atoi(match[4])
		if err != nil {
			return
		}
		if match[3] != "true" && match[3] != "false" {
			return markerSets, fmt.Errorf("unable to parse %v as bool\n", match[3])
		}

		markerSets = append(markerSets, MarkerSet{
			Name:        match[1],
			Label:       match[2],
			Hide:        match[3] == "true",
			Priority:    priority,
			DefaultIcon: match[5],
		})
	}

	return
}
