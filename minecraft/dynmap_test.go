/*
Filename: dynmap_test.go
Description: test cases for select dynmap functions
Created by: osh
        at: 11:32 on Sunday, the 15th of June, 2025.
Last edited 12:45 on Sunday, the 15th of June, 2025.
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
