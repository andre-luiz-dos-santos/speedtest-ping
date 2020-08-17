// Send a message to the BotFather to get the token
// Send a message to the Bot and get the chat ID using
// curl https://api.telegram.org/bot<TOKEN>/getUpdates
// /setprivacy must be disabled to receive group messages.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendMessage(text string) error {
	b, err := json.Marshal(map[string]interface{}{
		"chat_id": telegramChatID,
		"text":    text,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", telegramToken)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram sendMessage: %v", resp.Status)
	}

	return nil
}
