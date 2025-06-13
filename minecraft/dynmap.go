/*
Filename: dynmap.go
Description: interfaces with dynmap marker commands to allow for the
             creation, modification and deletion of dynmap markers
Created by: osh
        at: 15:39 on Friday, the 13th of June, 2025.
Last edited 17:30 on Friday, the 13th of June, 2025.
*/

package minecraft

type Marker struct {
	Label     string
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

// adds a marker given marker fields
func AddMarker(marker Marker) error

// builds the `/dmarker add` command, allowing for optional fields
func buildAddMarker(marker Marker) (string, error)

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
