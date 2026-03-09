package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Alert struct {
	Level   string // INFO, WARNING, CRITICAL
	Message string
	Time    time.Time
}

func SendDiscord(webhookURL string, alert Alert) error {
	color := 0x3498db // Blue
	if alert.Level == "CRITICAL" {
		color = 0xe74c3c // Red
	} else if alert.Level == "WARNING" {
		color = 0xf1c40f // Yellow
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       fmt.Sprintf("[%s] WireGate Alert", alert.Level),
				"description": alert.Message,
				"color":       color,
				"timestamp":   alert.Time.Format(time.RFC3339),
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord returned status: %d", resp.StatusCode)
	}

	return nil
}

func SendTelegram(token, chatID string, alert Alert) error {
	message := fmt.Sprintf("[%s] WireGate Alert\n\n%s\n\nTime: %s", 
		alert.Level, alert.Message, alert.Time.Format("2006-01-02 15:04:05"))

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram returned status: %d", resp.StatusCode)
	}

	return nil
}
