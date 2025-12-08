package main

import (
	"github.com/michaelfmnk/send-to-kindle-telegram-bot/bot"
	"log"
	"os"
	"strings"
)

func main() {
	// Check if SMTP insecure mode is enabled
	smtpInsecure := false
	if insecureEnv := os.Getenv("UBOT_SMTP_INSECURE"); insecureEnv != "" {
		smtpInsecure = strings.ToLower(insecureEnv) == "true" || insecureEnv == "1"
	}

	// Parse multi-Kindle configuration
	kindleDevices := parseKindleDevices(os.Getenv("UBOT_KINDLE_DEVICES"))

	unkindleBot := bot.SendToKindleBot{
		Token:         os.Getenv("UBOT_TELEGRAM_TOKEN"),
		EmailFrom:     os.Getenv("UBOT_EMAIL_FROM"),
		EmailTo:       os.Getenv("UBOT_EMAIL_TO"), // Fallback for single device
		KindleDevices: kindleDevices,                // Multi-device support
		SMTPHost:      os.Getenv("UBOT_SMTP_HOST"),
		SMTPPort:      os.Getenv("UBOT_SMTP_PORT"),
		Password:      os.Getenv("UBOT_PASSWORD"),
		SMTPInsecure:  smtpInsecure,
	}
	if err := unkindleBot.Start(); err != nil {
		log.Fatal("[ERROR] could not start telegram bot:", err)
	}
}

// parseKindleDevices parses UBOT_KINDLE_DEVICES into a map
// Format: "Device1:email1@kindle.com|Device2:email2@kindle.com"
// Example: "Kindle Paperwhite:user1@kindle.com|Kindle Oasis:user2@kindle.com"
func parseKindleDevices(devicesEnv string) map[string]string {
	devices := make(map[string]string)

	if devicesEnv == "" {
		return devices // Empty map, will use EmailTo as fallback
	}

	// Split by pipe separator
	devicePairs := strings.Split(devicesEnv, "|")
	for _, pair := range devicePairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// Split device name and email by colon
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			log.Printf("[WARN] Invalid device format (expected 'Name:email'): %s\n", pair)
			continue
		}

		deviceName := strings.TrimSpace(parts[0])
		deviceEmail := strings.TrimSpace(parts[1])

		if deviceName == "" || deviceEmail == "" {
			log.Printf("[WARN] Empty device name or email: %s\n", pair)
			continue
		}

		devices[deviceName] = deviceEmail
		log.Printf("[INFO] Registered Kindle device: %s -> %s\n", deviceName, deviceEmail)
	}

	return devices
}
