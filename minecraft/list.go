package minecraft

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(
	`^There are (\d+) of a max of (\d+) players online(?:\:\s*(.*))?\.?$`)

// Example:
// There are 3 of a max of 20 players online: OshDubh, MosEisley1976, Wizzeroo
//
// String is also stable
func extractInfo(s string) (string, string, []string) {
	m := re.FindStringSubmatch(s)
	if m == nil {
		return "", "", nil
	}

	var players []string
	if len(m[3]) != 0 {
		for _, p := range strings.Split(m[3], ", ") {
			p = strings.TrimSpace(p)
			if p != "" {
				players = append(players, p)
			}
		}
	}

	return m[1], m[2], players
}

func GetOnlinePlayerMsg(s string) string {
	numberOfPlayers, max, players := extractInfo(s)
	if numberOfPlayers == "" || numberOfPlayers == "0" {
		return s
	}

	for i, player := range players {
		players[i] = " - " + player
	}

	return fmt.Sprintf("There are %s of a max of %s players online:\n%s", numberOfPlayers, max, strings.Join(players, "\n"))
}
