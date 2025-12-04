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

	unkindleBot := bot.SendToKindleBot{
		Token:        os.Getenv("UBOT_TELEGRAM_TOKEN"),
		EmailFrom:    os.Getenv("UBOT_EMAIL_FROM"),
		EmailTo:      os.Getenv("UBOT_EMAIL_TO"),
		SMTPHost:     os.Getenv("UBOT_SMTP_HOST"),
		Password:     os.Getenv("UBOT_PASSWORD"),
		SMTPInsecure: smtpInsecure,
	}
	if err := unkindleBot.Start(); err != nil {
		log.Fatal("could not start telegram bot", err)
	}
}