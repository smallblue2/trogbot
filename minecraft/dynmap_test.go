/*
Filename: dynmap_test.go
Description: test cases for select dynmap functions
Created by: osh
        at: 11:32 on Sunday, the 15th of June, 2025.
Last edited 13:31 on Sunday, the 15th of June, 2025.
*/

package minecraft

import (
	"testing"
)

func TestParseWorlds(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wantE bool
		wants []MarkerWorld
	}{
		{
			name:  "actual input",
			input: `world world: loaded=true, enabled=true, title=overworld, center=112.0/77.0/528.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=trueworld DIM-1: loaded=true, enabled=true, title=the_nether, center=112.0/77.0/528.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=trueworld DIM1: loaded=true, enabled=true, title=the_end, center=112.0/77.0/528.0, extrazoomout=2, sendhealth=true, sendposition=true, protected=false, showborder=true`,
			wantE: false,
			wants: []MarkerWorld{
				{
					Name:         "world",
					Loaded:       true,
					Enabled:      true,
					Title:        "overworld",
					Center:       "112.0/77.0/528.0",
					ExtraZoomOut: 2,
					SendHealth:   true,
					SendPosition: true,
					Protected:    false,
					ShowBorder:   true,
				},
				{
					Name:         "DIM-1",
					Loaded:       true,
					Enabled:      true,
					Title:        "the_nether",
					Center:       "112.0/77.0/528.0",
					ExtraZoomOut: 2,
					SendHealth:   true,
					SendPosition: true,
					Protected:    false,
					ShowBorder:   true,
				},
				{
					Name:         "DIM1",
					Loaded:       true,
					Enabled:      true,
					Title:        "the_end",
					Center:       "112.0/77.0/528.0",
					ExtraZoomOut: 2,
					SendHealth:   true,
					SendPosition: true,
					Protected:    false,
					ShowBorder:   true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := parseMarkerWorlds(tc.input)

			if (err != nil) != tc.wantE {
				t.Errorf("\nerror: %vwanted error: %v", err, tc.wantE)
				return
			}

			if len(results) != len(tc.wants) {
				t.Errorf("\ngot: %v results\nwant %v", len(results), len(tc.wants))
				t.Logf("\nresults: %+v", results)
				return
			}

			for i, result := range results {
				if result != tc.wants[i] {
					t.Errorf("\ngot: %+v\nwant: %+v", result, tc.wants[i])
				}
			}
		})
	}
}

func TestParseMarkerIcons(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wantE bool
		wants []MarkerIcon
	}{
		{
			name:  "actual input",
			input: `anchor: label:"anchor", builtin:truebank: label:"bank", builtin:truebasket: label:"basket", builtin:truebed: label:"bed", builtin:truebeer: label:"beer", builtin:truebighouse: label:"bighouse", builtin:trueblueflag: label:"blueflag", builtin:truebomb: label:"bomb", builtin:truebookshelf: label:"bookshelf", builtin:truebricks: label:"bricks", builtin:truebronzemedal: label:"bronzemedal", builtin:truebronzestar: label:"bronzestar", builtin:truebuilding: label:"building", builtin:truecake: label:"cake", builtin:truecamera: label:"camera", builtin:truecart: label:"cart", builtin:truecaution: label:"caution", builtin:truechest: label:"chest", builtin:truechurch: label:"church", builtin:truecoins: label:"coins", builtin:truecomment: label:"comment", builtin:truecompass: label:"compass", builtin:trueconstruction: label:"construction", builtin:truecross: label:"cross", builtin:truecup: label:"cup", builtin:truecutlery: label:"cutlery", builtin:truedefault: label:"default", builtin:truediamond: label:"diamond", builtin:truedog: label:"dog", builtin:truedoor: label:"door", builtin:truedown: label:"down", builtin:truedrink: label:"drink", builtin:trueexclamation: label:"exclamation", builtin:truefactory: label:"factory", builtin:truefire: label:"fire", builtin:trueflower: label:"flower", builtin:truegear: label:"gear", builtin:truegoldmedal: label:"goldmedal", builtin:truegoldstar: label:"goldstar", builtin:truegreenflag: label:"greenflag", builtin:truehammer: label:"hammer", builtin:trueheart: label:"heart", builtin:truehouse: label:"house", builtin:truekey: label:"key", builtin:trueking: label:"king", builtin:trueleft: label:"left", builtin:truelightbulb: label:"lightbulb", builtin:truelighthouse: label:"lighthouse", builtin:truelock: label:"lock", builtin:trueminecart: label:"minecart", builtin:trueofflineuser: label:"offlineuser", builtin:trueorangeflag: label:"orangeflag", builtin:truepin: label:"pin", builtin:truepinkflag: label:"pinkflag", builtin:truepirateflag: label:"pirateflag", builtin:truepointdown: label:"pointdown", builtin:truepointleft: label:"pointleft", builtin:truepointright: label:"pointright", builtin:truepointup: label:"pointup", builtin:trueportal: label:"portal", builtin:truepurpleflag: label:"purpleflag", builtin:truequeen: label:"queen", builtin:trueredflag: label:"redflag", builtin:trueright: label:"right", builtin:trueruby: label:"ruby", builtin:truescales: label:"scales", builtin:trueshield: label:"shield", builtin:truesign: label:"sign", builtin:truesilvermedal: label:"silvermedal", builtin:truesilverstar: label:"silverstar", builtin:trueskull: label:"skull", builtin:truestar: label:"star", builtin:truesun: label:"sun", builtin:truetemple: label:"temple", builtin:truetheater: label:"theater", builtin:truetornado: label:"tornado", builtin:truetower: label:"tower", builtin:truetree: label:"tree", builtin:truetruck: label:"truck", builtin:trueup: label:"up", builtin:truewalk: label:"walk", builtin:truewarning: label:"warning", builtin:trueworld: label:"world", builtin:truewrench: label:"wrench", builtin:trueyellowflag: label:"yellowflag", builtin:true`,
			wantE: false,
			wants: []MarkerIcon{
				{Name: "anchor", Label: "anchor", Builtin: true},
				{Name: "bank", Label: "bank", Builtin: true},
				{Name: "basket", Label: "basket", Builtin: true},
				{Name: "bed", Label: "bed", Builtin: true},
				{Name: "beer", Label: "beer", Builtin: true},
				{Name: "bighouse", Label: "bighouse", Builtin: true},
				{Name: "blueflag", Label: "blueflag", Builtin: true},
				{Name: "bomb", Label: "bomb", Builtin: true},
				{Name: "bookshelf", Label: "bookshelf", Builtin: true},
				{Name: "bricks", Label: "bricks", Builtin: true},
				{Name: "bronzemedal", Label: "bronzemedal", Builtin: true},
				{Name: "bronzestar", Label: "bronzestar", Builtin: true},
				{Name: "building", Label: "building", Builtin: true},
				{Name: "cake", Label: "cake", Builtin: true},
				{Name: "camera", Label: "camera", Builtin: true},
				{Name: "cart", Label: "cart", Builtin: true},
				{Name: "caution", Label: "caution", Builtin: true},
				{Name: "chest", Label: "chest", Builtin: true},
				{Name: "church", Label: "church", Builtin: true},
				{Name: "coins", Label: "coins", Builtin: true},
				{Name: "comment", Label: "comment", Builtin: true},
				{Name: "compass", Label: "compass", Builtin: true},
				{Name: "construction", Label: "construction", Builtin: true},
				{Name: "cross", Label: "cross", Builtin: true},
				{Name: "cup", Label: "cup", Builtin: true},
				{Name: "cutlery", Label: "cutlery", Builtin: true},
				{Name: "default", Label: "default", Builtin: true},
				{Name: "diamond", Label: "diamond", Builtin: true},
				{Name: "dog", Label: "dog", Builtin: true},
				{Name: "door", Label: "door", Builtin: true},
				{Name: "down", Label: "down", Builtin: true},
				{Name: "drink", Label: "drink", Builtin: true},
				{Name: "exclamation", Label: "exclamation", Builtin: true},
				{Name: "factory", Label: "factory", Builtin: true},
				{Name: "fire", Label: "fire", Builtin: true},
				{Name: "flower", Label: "flower", Builtin: true},
				{Name: "gear", Label: "gear", Builtin: true},
				{Name: "goldmedal", Label: "goldmedal", Builtin: true},
				{Name: "goldstar", Label: "goldstar", Builtin: true},
				{Name: "greenflag", Label: "greenflag", Builtin: true},
				{Name: "hammer", Label: "hammer", Builtin: true},
				{Name: "heart", Label: "heart", Builtin: true},
				{Name: "house", Label: "house", Builtin: true},
				{Name: "key", Label: "key", Builtin: true},
				{Name: "king", Label: "king", Builtin: true},
				{Name: "left", Label: "left", Builtin: true},
				{Name: "lightbulb", Label: "lightbulb", Builtin: true},
				{Name: "lighthouse", Label: "lighthouse", Builtin: true},
				{Name: "lock", Label: "lock", Builtin: true},
				{Name: "minecart", Label: "minecart", Builtin: true},
				{Name: "offlineuser", Label: "offlineuser", Builtin: true},
				{Name: "orangeflag", Label: "orangeflag", Builtin: true},
				{Name: "pin", Label: "pin", Builtin: true},
				{Name: "pinkflag", Label: "pinkflag", Builtin: true},
				{Name: "pirateflag", Label: "pirateflag", Builtin: true},
				{Name: "pointdown", Label: "pointdown", Builtin: true},
				{Name: "pointleft", Label: "pointleft", Builtin: true},
				{Name: "pointright", Label: "pointright", Builtin: true},
				{Name: "pointup", Label: "pointup", Builtin: true},
				{Name: "portal", Label: "portal", Builtin: true},
				{Name: "purpleflag", Label: "purpleflag", Builtin: true},
				{Name: "queen", Label: "queen", Builtin: true},
				{Name: "redflag", Label: "redflag", Builtin: true},
				{Name: "right", Label: "right", Builtin: true},
				{Name: "ruby", Label: "ruby", Builtin: true},
				{Name: "scales", Label: "scales", Builtin: true},
				{Name: "shield", Label: "shield", Builtin: true},
				{Name: "sign", Label: "sign", Builtin: true},
				{Name: "silvermedal", Label: "silvermedal", Builtin: true},
				{Name: "silverstar", Label: "silverstar", Builtin: true},
				{Name: "skull", Label: "skull", Builtin: true},
				{Name: "star", Label: "star", Builtin: true},
				{Name: "sun", Label: "sun", Builtin: true},
				{Name: "temple", Label: "temple", Builtin: true},
				{Name: "theater", Label: "theater", Builtin: true},
				{Name: "tornado", Label: "tornado", Builtin: true},
				{Name: "tower", Label: "tower", Builtin: true},
				{Name: "tree", Label: "tree", Builtin: true},
				{Name: "truck", Label: "truck", Builtin: true},
				{Name: "up", Label: "up", Builtin: true},
				{Name: "walk", Label: "walk", Builtin: true},
				{Name: "warning", Label: "warning", Builtin: true},
				{Name: "world", Label: "world", Builtin: true},
				{Name: "wrench", Label: "wrench", Builtin: true},
				{Name: "yellowflag", Label: "yellowflag", Builtin: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := parseMarkerIcons(tc.input)

			if (err != nil) != tc.wantE {
				t.Errorf("\nerror: %vwanted error: %v", err, tc.wantE)
				return
			}

			if len(results) != len(tc.wants) {
				t.Errorf("\ngot: %v results\nwant %v", len(results), len(tc.wants))
				t.Logf("\nresults: %+v", results)
				return
			}

			for i, result := range results {
				if result != tc.wants[i] {
					t.Errorf("\ngot: %+v\nwant: %+v", result, tc.wants[i])
				}
			}
		})
	}
}

func TestParseMarkerSets(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wantE bool
		wants []MarkerSet
	}{
		{
			name:  "actual input",
			input: `markers: label:"Markers", hide:false, prio:0, deficon:default, persistent=true`,
			wantE: false,
			wants: []MarkerSet{
				{
					Name:        "markers",
					Label:       "Markers",
					Hide:        false,
					Priority:    0,
					DefaultIcon: "default",
					Persistent:  true,
				},
			},
		},
		{
			name:  "no persistent field",
			input: `markers: label:"Markers", hide:false, prio:0, deficon:default`,
			wantE: false,
			wants: []MarkerSet{
				{
					Name:        "markers",
					Label:       "Markers",
					Hide:        false,
					Priority:    0,
					DefaultIcon: "default",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := parseMarkerSets(tc.input)

			if (err != nil) != tc.wantE {
				t.Errorf("\nerror: %vwanted error: %v", err, tc.wantE)
				return
			}

			if len(results) != len(tc.wants) {
				t.Errorf("\ngot: %v results\nwant %v", len(results), len(tc.wants))
				t.Logf("\nresults: %+v", results)
				return
			}

			for i, result := range results {
				if result != tc.wants[i] {
					t.Errorf("\ngot: %+v\nwant: %+v", result, tc.wants[i])
				}
			}
		})
	}

}
