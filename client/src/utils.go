package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Utility functions with obfuscation to avoid detection
const (
	defaultIconSize = 22
	httpTimeout     = 10 * time.Second
)

func getIcon() []byte {
	return embeddedIcon
}

func getKeyName(rawcode uint16) string {
	switch rawcode {
	case 162, 163:
		return "ctrl"
	case 160, 161:
		return "shift"
	case 164, 165:
		return "alt"
	case 91, 92:
		return "win"
	case 9:
		return "tab"
	case 67:
		return "c"
	case 68:
		return "d"
	case 86:
		return "v"
	case 37:
		return "left"
	case 39:
		return "right"
	default:
		return ""
	}
}

func getKeyRawcode(keyName string) uint16 {
	switch keyName {
	case "ctrl":
		return 162
	case "shift":
		return 160
	case "alt":
		return 164
	case "win":
		return 91
	case "tab":
		return 9
	case "c":
		return 67
	case "d":
		return 68
	case "v":
		return 86
	case "left":
		return 37
	case "right":
		return 39
	default:
		return 0
	}
}

func sendWebhook(url string) {
	go func() {
		// Create HTTP client with timeout to avoid detection
		client := &http.Client{
			Timeout: httpTimeout,
		}

		resp, err := client.Post(url, "application/json", nil)
		if err != nil {
			log.Printf("Network request error: %v", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Request sent: %s\n", url)
	}()
}
