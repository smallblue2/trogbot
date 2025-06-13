package minecraft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/smallblue2/trogbot/config"
)

type Entry struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

var mu sync.RWMutex

func Load() ([]Entry, error) {
	mu.RLock()
	defer mu.RUnlock()

	b, err := os.ReadFile(config.WHITELIST_PATH)
	if err != nil {
		return nil, err
	}
	var list []Entry
	return list, json.Unmarshal(b, &list)
}

func Save(list []Entry) error {
	mu.Lock()
	defer mu.Unlock()

	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(config.WHITELIST_PATH, b, 0644)
}

func dashUUID(id32 string) string {
	if len(id32) != 32 {
		return id32
	}

	return id32[:8] + "-" +
		id32[8:12] + "-" +
		id32[12:16] + "-" +
		id32[16:20] + "-" +
		id32[20:]
}

func FetchUUID(ctx context.Context, name string) (string, error) {
	url := "https://api.minecraftservices.com/minecraft/profile/lookup/name/" + name
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("player not found")
	}

	var j struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return "", err
	}

	if len(j.ID) == 0 {
		return "", fmt.Errorf("player not found")
	}

	return dashUUID(j.ID), nil
}
