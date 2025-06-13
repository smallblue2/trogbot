package registry

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Definition() *discordgo.ApplicationCommand
	Run(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

var registry = struct {
	sync.RWMutex
	m map[string]Command
}{m: make(map[string]Command)}

func Register(c Command) {
	registry.Lock()
	name := c.Definition().Name
	registry.m[name] = c
	log.Printf("Registered command '%s'\n", name)
	registry.Unlock()
}

func Lookup(name string) (Command, bool) {
	registry.RLock()
	c, ok := registry.m[name]
	registry.RUnlock()
	return c, ok
}

func All() []Command {
	registry.RLock()
	out := make([]Command, 0, len(registry.m))
	for _, c := range registry.m {
		out = append(out, c)
	}
	registry.RUnlock()
	return out
}

func AllDefinitions() []*discordgo.ApplicationCommand {
	registry.RLock()
	out := make([]*discordgo.ApplicationCommand, 0, len(registry.m))
	for _, c := range registry.m {
		out = append(out, c.Definition())
	}
	registry.RUnlock()
	return out
}
