package minecraft

import (
	"log"

	"github.com/jltobler/go-rcon"
	"github.com/smallblue2/trogbot/config"
)

var client = rcon.NewClient(config.RCON_ADDR, config.RCON_PASS)

func Exec(cmd string) (string, error) {
	log.Printf("Running RCON command `%s`\n", cmd)
	return client.Send(cmd)
}
