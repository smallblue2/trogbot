package minecraft

import (
	"slices"
	"testing"
)

func TestExtractInfo(t *testing.T) {
	tests := []struct {
		name            string
		text            string
		expectedOnline  string
		expectedMax     string
		expectedPlayers []string
	}{
		{"None online", "There are 0 of a max of 20 players online.", "0", "20", nil},
		{"Some online", "There are 3 of a max of 5 players online: OshDubh, MosEisley1976, Wizzeroo", "3", "5", []string{"OshDubh", "MosEisley1976", "Wizzeroo"}},
		{"Wrong input", "This is wrong input", "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			online, max, players := extractInfo(tt.text)
			if online != tt.expectedOnline {
				t.Errorf("online = '%s'; expected = '%s'\n", online, tt.expectedOnline)
			}
			if max != tt.expectedMax {
				t.Errorf("max = '%s'; expected = '%s'\n", max, tt.expectedMax)
			}
			if !slices.Equal(players, tt.expectedPlayers) {
				t.Errorf("players = '%s'; expected = '%s'\n", players, tt.expectedPlayers)
			}
		})
	}
}
