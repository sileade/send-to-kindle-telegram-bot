package bot

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/scorredoira/email"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultSMTPPort = "587"
	tmpFilesPath    = "/files/"
)

var (
	// ErrNoToken - represents a validation error when Token not set
	ErrNoToken = errors.New("token for telegram bot not set")
	// ErrNoPassword - represents a validation error when Password not set
	ErrNoPassword = errors.New("password for email not set")
	// ErrNoEmailFrom - represents a validation error when EmailFrom not set
	ErrNoEmailFrom = errors.New("emailfrom not set")
	// ErrNoEmailTo - represents a validation error when EmailTo not set
	ErrNoEmailTo = errors.New("emailto not set")
	// ErrStartup - represents an error during bot initialization process
	ErrStartup = errors.New("could not create telebot instance")

	errConversion    = errors.New("could not convert file")
	supportedFormats = []string{"doc", "docx", "rtf", "htm", "html", "txt", "mobi", "pdf"}
)

// SendToKindleBot stores bot configuration
type SendToKindleBot struct {
	Token        string
	EmailFrom    string
	EmailTo      string
	SMTPHost     string
	SMTPPort     string
	Password     string
	SMTPInsecure bool
}

// Start starts bot. It is blocking.
// If there is an error during startup, returns it. Otherwise blocks
func (b *SendToKindleBot) Start() error {
	if err := b.verifyConfig(); err != nil {
		return err
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  b.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return ErrStartup
	}

	bot.Handle(tb.OnDocument, b.documentHandler(bot))
	bot.Start()

	return nil
}

func (b *SendToKindleBot) documentHandler(bot *tb.Bot) func(msg *tb.Message) {
	return func(msg *tb.Message) {
		doc := msg.Document
		nameParts := strings.Split(doc.FileName, ".")
		fileNameWithoutExtension := strings.Join(nameParts[:len(nameParts)-1], "")
		extension := nameParts[len(nameParts)-1]

		originalFilePath := tmpFilesPath + doc.FileName
		if err := bot.Download(&doc.File, originalFilePath); err != nil {
			log.Println("could not download file", err)
			respond(bot, msg, "Sorry. I could not download file")
			return
		}
		defer removeSilently(originalFilePath)

		fileToSend := originalFilePath
		if needToConvert(extension) {
			outputFilePath := tmpFilesPath + fileNameWithoutExtension + ".mobi"
			if err := convert(originalFilePath, outputFilePath); err != nil {
				log.Println("could not convert file", err)
				respond(bot, msg, "Sorry. I could not convert file")
				return
			}
			fileToSend = outputFilePath
			defer removeSilently(outputFilePath)
		}

		if err := b.sendFileViaEmail(fileToSend); err != nil {
			log.Println("could not send file", err)
			respond(bot, msg, "Sorry. I could not send file")
			return
		}
		respond(bot, msg, "File sent successfully!")
	}
}

func needToConvert(extension string) bool {
	for _, format := range supportedFormats {
		if format == extension {
			return false
		}
	}
	return true
}

func respond(bot *tb.Bot, m *tb.Message, text string) {
	if _, err := bot.Send(m.Sender, text); err != nil {
		log.Println(fmt.Sprintf("could not send a message to %d", m.Sender.ID), err)
	}
}

func convert(in, out string) error {
	cmd := exec.Command("ebook-convert", in, out)
	if err := cmd.Run(); err != nil {
		return err
	}
	if _, err := os.Stat(out); errors.Is(err, os.ErrNotExist) {
		return errConversion
	}
	return nil
}

func removeSilently(path string) {
	if err := os.Remove(path); err != nil {
		log.Println(fmt.Sprintf("could not delete file %s", path), err)
	}
}

func (b *SendToKindleBot) verifyConfig() error {
	if b.Token == "" {
		return ErrNoToken
	}
	if b.Password == "" {
		return ErrNoPassword
	}
	if b.EmailFrom == "" {
		return ErrNoEmailFrom
	}
	if b.EmailTo == "" {
		return ErrNoEmailTo
	}
	if b.SMTPPort == "" {
		b.SMTPPort = defaultSMTPPort
	}
	// Remove port from SMTPHost if it contains one
	if strings.Contains(b.SMTPHost, ":") {
		parts := strings.Split(b.SMTPHost, ":")
		b.SMTPHost = parts[0]
		if len(parts) > 1 && b.SMTPPort == defaultSMTPPort {
			b.SMTPPort = parts[1]
		}
	}
	return nil
}

func (b *SendToKindleBot) sendFileViaEmail(path string) error {
	msg := email.NewMessage("", "")
	msg.From = mail.Address{Name: "From", Address: b.EmailFrom}
	msg.To = []string{b.EmailTo}

	if err := msg.Attach(path); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", b.EmailFrom, b.Password, b.SMTPHost)
	addr := fmt.Sprintf("%s:%s", b.SMTPHost, b.SMTPPort)

	// Configure TLS
	tlsConfig := &tls.Config{
		ServerName:         b.SMTPHost,
		InsecureSkipVerify: b.SMTPInsecure,
	}

	// Send with custom TLS config
	if err := sendEmailWithTLS(addr, auth, msg, tlsConfig); err != nil {
		return err
	}
	return nil
}

// sendEmailWithTLS sends email with custom TLS configuration
func sendEmailWithTLS(addr string, auth smtp.Auth, msg *email.Message, tlsConfig *tls.Config) error {
	// This is a workaround since the email library doesn't expose TLS config
	// We need to use the standard approach
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(msg.From.Address); err != nil {
		return err
	}

	for _, to := range msg.To {
		if err = c.Rcpt(to); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg.Bytes())
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}