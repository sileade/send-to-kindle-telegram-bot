# 🚀 Quick Start Guide

## One-Time Setup (5 minutes)

### Step 1: Clone or Pull Latest Changes

```bash
cd /path/to/send-to-kindle-telegram-bot
git pull origin master
```

### Step 2: Create .env Configuration

```bash
cp .env.example .env
nano .env  # Edit with your settings
```

**Required fields to fill:**

| Field | Value | Example |
|-------|-------|----------|
| `UBOT_TELEGRAM_TOKEN` | Get from [@BotFather](https://t.me/botfather) | `123456789:ABCdefGHIjklMNOpqrsTUVwxyz` |
| `UBOT_EMAIL_FROM` | Your email address | `myemail@gmail.com` |
| `UBOT_PASSWORD` | **App-specific password** (not your main password!) | Use links in .env.example |
| `UBOT_EMAIL_TO` | Your Kindle email | `username@kindle.com` |
| `UBOT_SMTP_HOST` | SMTP server (without port) | `smtp.gmail.com` |
| `UBOT_SMTP_PORT` | Usually 587 | `587` |

### Step 3: Configure Kindle

1. Go to [Amazon Manage Your Content](https://www.amazon.com/gp/digital/fiona/)
2. Find "Approved Personal Document E-mail List"
3. Click "Add a new approved e-mail address"
4. Enter your `UBOT_EMAIL_FROM` address

### Step 4: Build and Start

```bash
# Build the container
docker compose build --no-cache

# Start the bot
docker compose up -d

# Check it's running
docker compose logs -f
```

You should see:
```
[INFO] Starting Send-to-Kindle bot...
[INFO] Using SMTP: smtp.gmail.com:587
[INFO] Bot successfully created and listening for documents...
```

## Testing (30 seconds)

1. Open Telegram and find your bot
2. Send a PDF file or any supported document
3. Bot should respond: `✅ File sent successfully to your Kindle!`
4. Check your Kindle device - the book should appear in a few minutes

## Ongoing Usage

### Check Bot Status

```bash
# View logs in real-time
docker compose logs -f

# Check container is running
docker ps | grep sendtokindle

# View resource usage
docker stats sendtokindle-bot
```

### Stop/Start Bot

```bash
# Stop
docker compose stop

# Start again
docker compose start

# Restart
docker compose restart
```

### Update to Latest Version

```bash
# Get latest code
git pull origin master

# Rebuild
docker compose build --no-cache

# Restart
docker compose restart
```

## Troubleshooting Quick Links

- **Book not being delivered?** → See [README.md Troubleshooting](README.md#troubleshooting)
- **Container won't start?** → Check logs: `docker compose logs sendtokindle`
- **Authentication failed?** → Use app-specific password, not main password
- **SMTP connection error?** → Check SMTP host (no port) and firewall
- **Need detailed guide?** → Read [REBUILD.md](REBUILD.md)

## Email Provider Quick Reference

### Gmail
- **SMTP Host:** `smtp.gmail.com`
- **Port:** `587`
- **Password:** [Create app password](https://myaccount.google.com/apppasswords)
- **Enable:** [Less secure app access](https://myaccount.google.com/lesssecureapps)

### Yandex
- **SMTP Host:** `smtp.yandex.com`
- **Port:** `587`
- **Password:** [Generate app password](https://passport.yandex.ru/)

### Outlook / Microsoft 365
- **SMTP Host:** `smtp.live.com` or `smtp.office365.com`
- **Port:** `587`
- **Password:** [Create app password](https://account.live.com/)

### Other Providers
Check your email provider's SMTP settings documentation.

## File Format Support

### Direct (no conversion needed)
✅ PDF, DOC, DOCX, RTF, HTML, HTM, TXT, MOBI

### Auto-converted to MOBI
✅ EPUB, FB2, AZW, and 100+ other formats (via Calibre)

## Security Notes

⚠️ **Never commit your .env file to git!** It contains secrets.

The `.gitignore` should already exclude it, but verify:
```bash
git status | grep .env  # Should not appear
```

## Getting Help

1. Check the [Troubleshooting Guide](README.md#troubleshooting)
2. View detailed logs: `docker compose logs sendtokindle --tail 50`
3. Check [CHANGELOG.md](CHANGELOG.md) for recent fixes
4. Review [REBUILD.md](REBUILD.md) for deployment details

---

**That's it!** Your bot should now be running and ready to send books to your Kindle! 📚
