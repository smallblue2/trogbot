/*
Filename: dynmap.go
Description: interfaces with dynmap marker commands to allow for the
             creation, modification and deletion of dynmap markers
Created by: osh
        at: 15:39 on Friday, the 13th of June, 2025.
Last edited 19:33 on Friday, the 13th of June, 2025.
*/

package minecraft

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
	Label     string
	LabelName string // unsure if label and LabelName are different, it might just be the same field listed twice
	Builtin   bool
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

// adds a marker given marker fields
func AddMarker(marker Marker) error

// builds the `/dmarker add` command, allowing for optional fields
func buildAddMarker(marker Marker) (string, error)

// retrieves the currently available worlds using `/dmap worldlist` to populat ethe slash command options
// looks like (newline seperated):
//
// world <world id>: loaded=<bool>, enabled=<bool>, title=<the dimension name>, center=<x.x/y.y/z.z>, extrazoomout=<int>, sendpositon=<bool>, protected=<bool>, showborder=<bool>
// eg:
// world world_172429123: loaded=true, enabled=true, title=world, center=-32.0/70.0/-80.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=true
func GetMarkerWorlds() ([]MarkerWorld, error)

// retrieves currently available marker icons using `/dmarker icons` to populate the slash command options
// looks like (newline seperated):
//
// <label>: label:"<label>", builtin:<bool>
// eg:
// anchor: label:"anchor", builtin:true
func GetMarkerIcons() ([]MarkerIcon, error)

// retrieves currently available marker sets using `/dmarker listsets` to populate the slash command options
// looks like (newline seperated):
//
// <marker set name>: label:"<label display name>", hide:<bool>, prio:<int>, deficon:<default icon name>
// eg:
// markers: label:"Markers", hide:false, prio:0, deficon:default
func GetMarkerSets() ([]MarkerSet, error)
