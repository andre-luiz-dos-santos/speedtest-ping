package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var reqJSON map[string]interface{}

	// Read JSON request.
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
	err = dec.Decode(&reqJSON)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get IP address from addr (IP:port).
	addr, _, err := net.SplitHostPort(reqJSON["addr"].(string))
	if err != nil {
		http.Error(w, "Invalid remote address", http.StatusBadRequest)
		return
	}

	// Get information about addr from IXC.
	info, err := getInfoForAddr(addr)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "getInfoForAddr: %v\n", err)
		return
	}

	// Add information from the speedtest site.
	if _, ok := reqJSON["start"]; ok {
		info += "--\n"
		info += "Iniciado\n"
	}
	if _, ok := reqJSON["download"]; ok {
		info += "--\n"
		info += fmt.Sprintf("Down: %v\n", reqJSON["download"])
		info += fmt.Sprintf("Up: %v\n", reqJSON["upload"])
		info += fmt.Sprintf("Ping: %v\n", reqJSON["ping"])
		info += fmt.Sprintf("Jitter: %v\n", reqJSON["jitter"])
	}

	// Send everything above to Telegram.
	err = sendMessage(info)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "sendMessage: %v\n", err)
		return
	}
}
