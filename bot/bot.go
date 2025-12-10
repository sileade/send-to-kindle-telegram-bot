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
	"sync"
	"time"
)

const (
	defaultSMTPPort       = "587"
	defaultTmpFilesPath   = "/files/"
	buttonsPerRow         = 2
	callbackDataPrefix    = "send_kindle:"
	maxFileNameLength     = 255
	maxDeviceNameLength   = 100
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
	// ErrInvalidFileName - represents an error for invalid filename
	ErrInvalidFileName = errors.New("invalid filename")

	errConversion = errors.New("could not convert file")
	// FIXED: Changed from MOBI to EPUB as Amazon discontinued MOBI support in Send to Kindle service
	// EPUB is the recommended format for modern Kindle devices (including Paperwhite 2024)
	supportedFormats = []string{"epub", "doc", "docx", "rtf", "htm", "html", "txt", "pdf"}
)

// SendToKindleBot stores bot configuration
type SendToKindleBot struct {
	Token          string
	EmailFrom      string
	EmailTo        string            // Single device (fallback)
	KindleDevices  map[string]string // Multiple devices: name -> email
	SMTPHost       string
	SMTPPort       string
	Password       string
	SMTPInsecure   bool
	bot            *tb.Bot
	fileStateCache map[int]map[string]string // userID -> {filePath, originalFileName}
	cacheMutex     sync.RWMutex              // FIXED: Added mutex for thread-safe access
	tmpFilesPath   string                    // FIXED: Made configurable
}

// Start starts bot. It is blocking.
// If there is an error during startup, returns it. Otherwise blocks
func (b *SendToKindleBot) Start() error {
	if err := b.verifyConfig(); err != nil {
		return err
	}

	// FIXED: Set default tmp files path if not set
	if b.tmpFilesPath == "" {
		b.tmpFilesPath = defaultTmpFilesPath
	}

	log.Println("[INFO] Starting Send-to-Kindle bot...")
	log.Printf("[INFO] Using SMTP: %s:%s\n", b.SMTPHost, b.SMTPPort)
	log.Printf("[INFO] Using temporary files path: %s\n", b.tmpFilesPath)

	// FIXED: Warn if insecure TLS is enabled
	if b.SMTPInsecure {
		log.Println("[WARN] SMTP insecure mode is enabled - TLS certificate verification is disabled!")
	}

	// Initialize file state cache
	b.fileStateCache = make(map[int]map[string]string)

	// Log available Kindle devices
	if len(b.KindleDevices) > 0 {
		log.Printf("[INFO] Available Kindle devices: %d\n", len(b.KindleDevices))
		for name := range b.KindleDevices {
			log.Printf("[DEBUG]   - %s\n", name)
		}
	} else if b.EmailTo != "" {
		// FIXED: Mask email in logs for security
		log.Printf("[INFO] Using single Kindle device: %s\n", maskEmail(b.EmailTo))
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  b.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return ErrStartup
	}
	b.bot = bot

	log.Println("[INFO] Bot successfully created and listening for documents...")
	bot.Handle(tb.OnDocument, b.documentHandler(bot))
	// Handle callback queries for device selection
	bot.Handle(tb.OnCallback, b.callbackHandler(bot))
	bot.Start()

	return nil
}

func (b *SendToKindleBot) documentHandler(bot *tb.Bot) func(msg *tb.Message) {
	return func(msg *tb.Message) {
		doc := msg.Document
		userID := msg.Sender.ID
		log.Printf("[DEBUG] Received document: %s from user %d\n", doc.FileName, userID)

		// FIXED: Validate and sanitize filename
		sanitizedFileName, err := sanitizeFileName(doc.FileName)
		if err != nil {
			log.Printf("[ERROR] Invalid filename: %v\n", err)
			respond(bot, msg, "❌ Invalid filename. Please check the file and try again.")
			return
		}

		// Get file extension and normalize to lowercase
		extension := strings.ToLower(filepath.Ext(sanitizedFileName))
		// Remove the leading dot if present
		if strings.HasPrefix(extension, ".") {
			extension = extension[1:]
		}

		// Get filename without extension
		fileNameWithoutExtension := strings.TrimSuffix(sanitizedFileName, filepath.Ext(sanitizedFileName))

		// Ensure tmpFilesPath exists
		if err := ensureDirectory(b.tmpFilesPath); err != nil {
			log.Printf("[ERROR] Could not create directory %s: %v\n", b.tmpFilesPath, err)
			respond(bot, msg, "❌ System error: could not prepare file storage")
			return
		}

		originalFilePath := filepath.Join(b.tmpFilesPath, sanitizedFileName)
		if err := bot.Download(&doc.File, originalFilePath); err != nil {
			log.Printf("[ERROR] Could not download file: %v\n", err)
			respond(bot, msg, "❌ Could not download file")
			return
		}

		fileToSend := originalFilePath
		if needToConvert(extension) {
			// FIXED: Changed from MOBI to EPUB format
			log.Printf("[DEBUG] Converting %s to EPUB format...\n", extension)
			outputFilePath := filepath.Join(b.tmpFilesPath, fileNameWithoutExtension+".epub")
			if err := convert(originalFilePath, outputFilePath); err != nil {
				log.Printf("[ERROR] Could not convert file: %v\n", err)
				respond(bot, msg, "❌ Could not convert file")
				removeSilently(originalFilePath)
				return
			}
			fileToSend = outputFilePath
		}

		// Store file info for callback handler (FIXED: with mutex)
		b.cacheMutex.Lock()
		if _, exists := b.fileStateCache[userID]; !exists {
			b.fileStateCache[userID] = make(map[string]string)
		}
		b.fileStateCache[userID]["filePath"] = fileToSend
		b.fileStateCache[userID]["originalFileName"] = sanitizedFileName
		b.fileStateCache[userID]["originalFilePath"] = originalFilePath
		b.cacheMutex.Unlock()

		// If only one device, send directly
		if len(b.KindleDevices) <= 1 && b.EmailTo != "" {
			log.Printf("[DEBUG] Sending file via email to %s...\n", maskEmail(b.EmailTo))
			if err := b.sendToKindle(fileToSend, sanitizedFileName, b.EmailTo); err != nil {
				log.Printf("[ERROR] Could not send file: %v\n", err)
				respond(bot, msg, "❌ Could not send file. Check logs for details")
				b.cleanupFiles(userID)
				return
			}
			respond(bot, msg, "✅ File sent successfully to your Kindle!")
			log.Printf("[INFO] Successfully sent %s to %s\n", sanitizedFileName, maskEmail(b.EmailTo))
			b.cleanupFiles(userID)
			return
		}

		// If multiple devices, show selection buttons
		if len(b.KindleDevices) > 1 {
			b.showDeviceSelection(bot, msg, userID)
			return
		}

		// No devices configured
		respond(bot, msg, "❌ No Kindle devices configured")
		b.cleanupFiles(userID)
	}
}

func (b *SendToKindleBot) showDeviceSelection(bot *tb.Bot, msg *tb.Message, userID int) {
	var buttons []tb.InlineButton

	for deviceName, deviceEmail := range b.KindleDevices {
		// Create callback data: "send_kindle:deviceName"
		callbackData := fmt.Sprintf("%s%s", callbackDataPrefix, deviceName)
		button := tb.InlineButton{
			Text: deviceName,
			Data: callbackData,
		}
		buttons = append(buttons, button)
		log.Printf("[DEBUG] Device button: %s (%s)\n", deviceName, maskEmail(deviceEmail))
	}

	// FIXED: Proper button layout with configurable buttons per row
	inlineKeys := make([][]tb.InlineButton, 0)
	for i := 0; i < len(buttons); i += buttonsPerRow {
		end := i + buttonsPerRow
		if end > len(buttons) {
			end = len(buttons)
		}
		inlineKeys = append(inlineKeys, buttons[i:end])
	}

	inlineMarkup := &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	}

	b.cacheMutex.RLock()
	originalFileName := b.fileStateCache[userID]["originalFileName"]
	b.cacheMutex.RUnlock()

	responseMsg := fmt.Sprintf("📱 Which Kindle device would you like to send '%s' to?\n\nSelect one:",
		originalFileName)
	if _, err := bot.Send(msg.Sender, responseMsg, inlineMarkup); err != nil {
		log.Printf("[ERROR] Could not send device selection: %v\n", err)
		respond(bot, msg, "❌ Could not show device selection. Please try again.")
	}
}

func (b *SendToKindleBot) callbackHandler(bot *tb.Bot) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		userID := c.Sender.ID
		callbackData := c.Data

		log.Printf("[DEBUG] Callback from user %d: %s\n", userID, callbackData)

		if !strings.HasPrefix(callbackData, callbackDataPrefix) {
			log.Printf("[DEBUG] Unknown callback: %s\n", callbackData)
			return
		}

		deviceName := strings.TrimPrefix(callbackData, callbackDataPrefix)
		deviceEmail, exists := b.KindleDevices[deviceName]
		if !exists {
			log.Printf("[ERROR] Device not found: %s\n", deviceName)
			bot.Respond(c, &tb.CallbackResponse{})
			bot.Send(c.Sender, "❌ Device not found")
			return
		}

		// Get file info from cache (FIXED: with mutex)
		b.cacheMutex.RLock()
		fileInfo, exists := b.fileStateCache[userID]
		b.cacheMutex.RUnlock()

		if !exists {
			log.Printf("[ERROR] No file in cache for user %d\n", userID)
			bot.Respond(c, &tb.CallbackResponse{})
			bot.Send(c.Sender, "❌ File not found. Please send it again.")
			return
		}

		filePath := fileInfo["filePath"]
		originalFileName := fileInfo["originalFileName"]

		// Send to selected device
		log.Printf("[DEBUG] Sending file to %s (%s)...\n", deviceName, maskEmail(deviceEmail))
		if err := b.sendToKindle(filePath, originalFileName, deviceEmail); err != nil {
			log.Printf("[ERROR] Could not send file to %s: %v\n", deviceName, err)
			bot.Respond(c, &tb.CallbackResponse{})
			bot.Send(c.Sender, fmt.Sprintf("❌ Could not send to %s. Try again.", deviceName))
			return
		}

		// Notify success
		bot.Respond(c, &tb.CallbackResponse{})
		bot.Send(c.Sender, fmt.Sprintf("✅ Book sent to %s!", deviceName))
		log.Printf("[INFO] Successfully sent %s to %s (%s)\n", originalFileName, deviceName, maskEmail(deviceEmail))

		// Cleanup
		b.cleanupFiles(userID)
	}
}

func (b *SendToKindleBot) sendToKindle(filePath string, originalFileName string, kindleEmail string) error {
	log.Printf("[DEBUG] Sending file via email to %s...\n", maskEmail(kindleEmail))

	// Create email with proper subject line
	subject := fmt.Sprintf("Book: %s", originalFileName)
	msg := email.NewMessage(subject, "")
	msg.From = mail.Address{Name: "Send-to-Kindle Bot", Address: b.EmailFrom}
	msg.To = []string{kindleEmail}

	if err := msg.Attach(filePath); err != nil {
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

func (b *SendToKindleBot) cleanupFiles(userID int) {
	b.cacheMutex.Lock()
	defer b.cacheMutex.Unlock()

	if fileInfo, exists := b.fileStateCache[userID]; exists {
		if filePath, ok := fileInfo["filePath"]; ok {
			removeSilently(filePath)
		}
		if origPath, ok := fileInfo["originalFilePath"]; ok {
			removeSilently(origPath)
		}
		delete(b.fileStateCache, userID)
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
		log.Printf("[ERROR] Could not send message to user %d: %v\n", m.Sender.ID, err)
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
	// At least one destination email required
	if b.EmailTo == "" && len(b.KindleDevices) == 0 {
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

// FIXED: Added sanitizeFileName to prevent path traversal attacks
func sanitizeFileName(fileName string) (string, error) {
	if fileName == "" {
		return "", ErrInvalidFileName
	}

	// Check length
	if len(fileName) > maxFileNameLength {
		return "", ErrInvalidFileName
	}

	// Remove path separators to prevent directory traversal
	fileName = filepath.Base(fileName)

	// Remove any remaining dangerous characters
	fileName = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1 // Remove control characters
		}
		switch r {
		case '<', '>', ':', '"', '|', '?', '*', 0:
			return -1 // Remove invalid filename characters
		}
		return r
	}, fileName)

	if fileName == "" {
		return "", ErrInvalidFileName
	}

	return fileName, nil
}

// FIXED: Added maskEmail to hide sensitive information in logs
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	// Show only domain part
	return "***@" + parts[1]
}
