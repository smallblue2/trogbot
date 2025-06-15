package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/smallblue2/trogbot/config"
)

type BackupNotification struct {
	Status    string `json:"status"`
	TimeUntil int    `json:"time_until,omitempty"`
}

func sendServerMessage(msg string) {
	if config.DISCORD_SESSION == nil || config.BOT_CHANNEL_ID == "" {
		log.Println("Cannot send message: DISCORD_SESSION or BOT_CHANNEL_ID not initialised")
		return
	}

	_, err := config.DISCORD_SESSION.ChannelMessageSend(config.BOT_CHANNEL_ID, msg)
	if err != nil {
		log.Printf("Failed to send message to Discord: %s\n", err)
	}
	log.Printf("Send message to Discord:\n'%s'\n", msg)
}

func backupNotificationHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method, expecting POST", http.StatusMethodNotAllowed)
		return
	}

	var notif BackupNotification
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	switch notif.Status {
	case "starting":
		sendServerMessage(fmt.Sprintf("⚠️ Backup starting in %d seconds! Please find a safe spot.", notif.TimeUntil))
	case "started":
		sendServerMessage("🚧 Backup started. Server saving data.")
	case "finished":
		sendServerMessage("✅ Backup finished! Server starting again.")
	default:
		http.Error(w, "Unknown status", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)

}

func StartServer() {
	http.HandleFunc("/backup", backupNotificationHandle)
}
