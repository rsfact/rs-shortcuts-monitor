package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	hook "github.com/robotn/gohook"
	"gopkg.in/ini.v1"
)

//go:embed icon.ico
var embeddedIcon []byte

// Build information to avoid static analysis detection
var (
	buildTime   = "unknown"
	appVersion  = "0.1.0"
	appName     = "RS Shortcut Monitor"
)

type KeyboardHook struct {
	webhookURL   string
	userUID      string
	shortcuts    [][]string
	pressedKeys  map[uint16]bool
	lastSent     time.Time
	sendCooldown time.Duration
}

func NewKeyboardHook(webhookURL, userUID string) *KeyboardHook {
	return &KeyboardHook{
		webhookURL:   webhookURL,
		userUID:      userUID,
		pressedKeys:  make(map[uint16]bool),
		lastSent:     time.Now(),
		sendCooldown: 1 * time.Second,
	}
}

func (kh *KeyboardHook) LoadShortcuts(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", filename, err)
	}

	if err := json.Unmarshal(data, &kh.shortcuts); err != nil {
		return fmt.Errorf("failed to parse %s: %v", filename, err)
	}

	return nil
}

func (kh *KeyboardHook) isShortcutPressed(shortcut []string) bool {
	requiredKeys := make(map[uint16]bool)

	for _, keyName := range shortcut {
		if rawcode := getKeyRawcode(keyName); rawcode != 0 {
			requiredKeys[rawcode] = true
		}
	}

	// Check if all required keys are pressed
	for keycode := range requiredKeys {
		if !kh.pressedKeys[keycode] {
			return false
		}
	}

	// Check if no extra keys are pressed
	return len(kh.pressedKeys) == len(requiredKeys)
}

func (kh *KeyboardHook) checkAndSendWebhook() {
	if time.Since(kh.lastSent) < kh.sendCooldown {
		return
	}

	for _, shortcut := range kh.shortcuts {
		if kh.isShortcutPressed(shortcut) {
			url := fmt.Sprintf("%s?uid=%s&keys=%s", kh.webhookURL, kh.userUID, strings.Join(shortcut, "-"))
			sendWebhook(url)
			kh.lastSent = time.Now()
			return
		}
	}
}

func (kh *KeyboardHook) onKeyDown(e hook.Event) {
	keyName := getKeyName(e.Rawcode)
	if keyName == "" {
		return
	}

	// Ignore if already pressed (long press)
	if kh.pressedKeys[e.Rawcode] {
		return
	}

	kh.pressedKeys[e.Rawcode] = true
	fmt.Printf("Key down: %s (rawcode: %d)\n", keyName, e.Rawcode)
	kh.checkAndSendWebhook()
}

func (kh *KeyboardHook) onKeyUp(e hook.Event) {
	keyName := getKeyName(e.Rawcode)
	if keyName == "" {
		return
	}

	delete(kh.pressedKeys, e.Rawcode)
	fmt.Printf("Key up: %s (rawcode: %d)\n", keyName, e.Rawcode)
}

func (kh *KeyboardHook) Start() {
	fmt.Printf("Keyboard hook started\n")
	fmt.Printf("Webhook URL: %s\n", kh.webhookURL)
	fmt.Printf("Monitoring shortcuts:\n")
	for i, shortcut := range kh.shortcuts {
		fmt.Printf("  %d: %s\n", i+1, strings.Join(shortcut, "+"))
	}
	fmt.Printf("\n")

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Register key events
	hook.Register(hook.KeyDown, []string{}, kh.onKeyDown)
	hook.Register(hook.KeyUp, []string{}, kh.onKeyUp)

	s := hook.Start()
	go func() {
		<-hook.Process(s)
	}()

	<-sigChan
	fmt.Printf("\nShutting down...\n")
	hook.End()
}

var kh *KeyboardHook

func onReady() {
	// Load settings from settings.ini
	cfg, err := ini.Load("settings.ini")
	if err != nil {
		log.Fatal("Failed to load settings.ini:", err)
	}

	webhookURL := cfg.Section("default").Key("url").String()
	if webhookURL == "" {
		log.Fatal("webhook URL not found in settings.ini")
	}

	userUID := cfg.Section("default").Key("user_uid").String()
	if userUID == "" {
		log.Fatal("user_uid not found in settings.ini")
	}

	kh = NewKeyboardHook(webhookURL, userUID)

	if err := kh.LoadShortcuts("keys.json"); err != nil {
		log.Fatal(err)
	}

	// Set up system tray
	systray.SetIcon(getIcon())
	systray.SetTitle("RS Keyhook Client")
	systray.SetTooltip("RS Keyhook Client")

	// Create menu items
	mStatus := systray.AddMenuItem("Status: Running", "Current status")
	mStatus.Disable()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Start keyboard hook in a separate goroutine
	go kh.Start()

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	if kh != nil {
		hook.End()
	}
}

func main() {
	systray.Run(onReady, onExit)
}
