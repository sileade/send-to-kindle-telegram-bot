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
	"path/filepath"
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
	// ErrNoSMTPHost - represents a validation error when SMTPHost not set
	ErrNoSMTPHost = errors.New("smtp host not set")
	// ErrStartup - represents an error during bot initialization process
	ErrStartup = errors.New("could not create telebot instance")

	errConversion    = errors.New("could not convert file")
	// FIXED: Changed from MOBI to EPUB as Amazon discontinued MOBI support in Send to Kindle service
	// EPUB is the recommended format for modern Kindle devices (including Paperwhite 2024)
	supportedFormats = []string{"epub", "doc", "docx", "rtf", "htm", "html", "txt", "pdf"}
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

	log.Println("[INFO] Starting Send-to-Kindle bot...")
	log.Printf("[INFO] Using SMTP: %s:%s\n", b.SMTPHost, b.SMTPPort)

	bot, err := tb.NewBot(tb.Settings{
		Token:  b.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return ErrStartup
	}

	log.Println("[INFO] Bot successfully created and listening for documents...")
	bot.Handle(tb.OnDocument, b.documentHandler(bot))
	bot.Start()

	return nil
}

func (b *SendToKindleBot) documentHandler(bot *tb.Bot) func(msg *tb.Message) {
	return func(msg *tb.Message) {
		doc := msg.Document
		log.Printf("[DEBUG] Received document: %s (size: %d bytes)\n", doc.FileName, doc.Size)

		// Get file extension and normalize to lowercase
		extension := strings.ToLower(filepath.Ext(doc.FileName))
		// Remove the leading dot if present
		if strings.HasPrefix(extension, ".") {
			extension = extension[1:]
		}

		// Get filename without extension
		fileNameWithoutExtension := strings.TrimSuffix(doc.FileName, filepath.Ext(doc.FileName))

		// Ensure tmpFilesPath exists
		if err := ensureDirectory(tmpFilesPath); err != nil {
			log.Printf("[ERROR] Could not create directory %s: %v\n", tmpFilesPath, err)
			respond(bot, msg, "Sorry. System error: could not prepare file storage")
			return
		}

		originalFilePath := filepath.Join(tmpFilesPath, doc.FileName)
		if err := bot.Download(&doc.File, originalFilePath); err != nil {
			log.Printf("[ERROR] Could not download file: %v\n", err)
			respond(bot, msg, "Sorry. I could not download file")
			return
		}
		defer removeSilently(originalFilePath)

		fileToSend := originalFilePath
		if needToConvert(extension) {
			// FIXED: Changed from MOBI to EPUB format
			log.Printf("[DEBUG] Converting %s to EPUB format...\n", extension)
			outputFilePath := filepath.Join(tmpFilesPath, fileNameWithoutExtension+".epub")
			if err := convert(originalFilePath, outputFilePath); err != nil {
				log.Printf("[ERROR] Could not convert file: %v\n", err)
				respond(bot, msg, "Sorry. I could not convert file")
				return
			}
			fileToSend = outputFilePath
			defer removeSilently(outputFilePath)
		}

		log.Printf("[DEBUG] Sending file via email to %s...\n", b.EmailTo)
		if err := b.sendFileViaEmail(fileToSend, doc.FileName); err != nil {
			log.Printf("[ERROR] Could not send file: %v\n", err)
			respond(bot, msg, "Sorry. I could not send file. Check logs for details")
			return
		}
		respond(bot, msg, "✅ File sent successfully to your Kindle!")
		log.Printf("[INFO] Successfully sent %s to %s\n", doc.FileName, b.EmailTo)
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
		log.Printf("[ERROR] Could not send message to %d: %v\n", m.Sender.ID, err)
	}
}

func convert(in, out string) error {
	log.Printf("[DEBUG] Running ebook-convert: %s -> %s\n", in, out)
	cmd := exec.Command("ebook-convert", in, out)
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] ebook-convert error: %v\n", err)
		return err
	}
	if _, err := os.Stat(out); errors.Is(err, os.ErrNotExist) {
		log.Printf("[ERROR] Conversion failed: output file not created\n")
		return errConversion
	}
	return nil
}

func removeSilently(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("[WARN] Could not delete file %s: %v\n", path, err)
	}
}

func ensureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
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
	if b.SMTPHost == "" {
		return ErrNoSMTPHost
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

func (b *SendToKindleBot) sendFileViaEmail(path string, originalFileName string) error {
	// Create email with proper subject line
	subject := fmt.Sprintf("Book: %s", originalFileName)
	msg := email.NewMessage(subject, "")
	msg.From = mail.Address{Name: "Send-to-Kindle Bot", Address: b.EmailFrom}
	msg.To = []string{b.EmailTo}

	if err := msg.Attach(path); err != nil {
		log.Printf("[ERROR] Could not attach file: %v\n", err)
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
	// Dial to SMTP server
	c, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("[ERROR] Could not connect to SMTP server %s: %v\n", addr, err)
		return fmt.Errorf("could not connect to SMTP server: %w", err)
	}
	defer c.Close()

	// Start TLS connection
	if err = c.StartTLS(tlsConfig); err != nil {
		log.Printf("[ERROR] Could not start TLS: %v\n", err)
		return fmt.Errorf("could not start TLS: %w", err)
	}

	// Authenticate after TLS
	if err = c.Auth(auth); err != nil {
		log.Printf("[ERROR] Authentication failed: %v\n", err)
		return fmt.Errorf("authentication failed (check email and password): %w", err)
	}

	// Send mail
	if err = c.Mail(msg.From.Address); err != nil {
		log.Printf("[ERROR] Could not set sender: %v\n", err)
		return fmt.Errorf("could not set sender: %w", err)
	}

	// Add recipients
	for _, to := range msg.To {
		if err = c.Rcpt(to); err != nil {
			log.Printf("[ERROR] Could not add recipient %s: %v\n", to, err)
			return fmt.Errorf("could not add recipient: %w", err)
		}
	}

	// Send data
	w, err := c.Data()
	if err != nil {
		log.Printf("[ERROR] Could not start data transmission: %v\n", err)
		return fmt.Errorf("could not start data transmission: %w", err)
	}

	_, err = w.Write(msg.Bytes())
	if err != nil {
		log.Printf("[ERROR] Could not write message data: %v\n", err)
		return fmt.Errorf("could not write message data: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[ERROR] Could not close data transmission: %v\n", err)
		return fmt.Errorf("could not close data transmission: %w", err)
	}

	// Quit
	if err = c.Quit(); err != nil {
		log.Printf("[ERROR] Could not close SMTP connection: %v\n", err)
		return fmt.Errorf("could not close SMTP connection: %w", err)
	}

	log.Printf("[DEBUG] Email sent successfully\n")
	return nil
}
