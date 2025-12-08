# Rebuilding and Deploying the Docker Container

## Quick Start

Follow these steps to rebuild the container with the latest fixes:

### 1. Pull Latest Changes

```bash
cd /path/to/send-to-kindle-telegram-bot
git pull origin master
```

### 2. Setup Environment File

If you don't have a `.env` file yet:

```bash
cp .env.example .env
```

Then edit `.env` with your actual configuration:

```bash
nano .env  # or use your favorite editor
```

Make sure to fill in:
- `UBOT_TELEGRAM_TOKEN` - Your Telegram bot token
- `UBOT_EMAIL_FROM` - Your email address
- `UBOT_PASSWORD` - Your app-specific password
- `UBOT_EMAIL_TO` - Your Kindle email address
- `UBOT_SMTP_HOST` - Your SMTP server
- `UBOT_SMTP_PORT` - SMTP port (usually 587)

### 3. Rebuild Docker Container

#### Option A: Using Docker Compose (Recommended)

```bash
# Clean build (removes old images)
docker compose build --no-cache

# Start the bot
docker compose up -d

# Check logs
docker compose logs -f
```

#### Option B: Using Docker Directly

```bash
# Build the image
docker build -t sendtokindle:latest .

# Stop existing container (if running)
docker stop sendtokindle-bot || true
docker rm sendtokindle-bot || true

# Run the new container
docker run -d \
  --name sendtokindle-bot \
  --restart unless-stopped \
  --env-file .env \
  -v "$(pwd)/files:/files" \
  sendtokindle:latest

# Check logs
docker logs -f sendtokindle-bot
```

## Verifying the Installation

### 1. Check Container Status

```bash
docker ps | grep sendtokindle
```

You should see the container running:
```
CONTAINER ID   IMAGE                     STATUS
abc123def     sendtokindle:latest       Up 2 minutes
```

### 2. Check Logs for Errors

```bash
docker compose logs sendtokindle
```

Look for these success messages:
```
[INFO] Starting Send-to-Kindle bot...
[INFO] Using SMTP: smtp.gmail.com:587
[INFO] Bot successfully created and listening for documents...
```

### 3. Test the Bot

1. Open Telegram and find your bot
2. Send a test document (PDF, EPUB, or other supported format)
3. Check the logs for debug messages:
   ```bash
   docker compose logs -f sendtokindle
   ```
4. The bot should respond with a confirmation message

## Troubleshooting

### Container Won't Start

```bash
# Check detailed logs
docker compose logs sendtokindle

# Rebuild from scratch
docker compose down
docker compose build --no-cache
docker compose up
```

### Authentication Failed

```
[ERROR] Authentication failed (check email and password)
```

Solutions:
- Use an app-specific password, not your main email password
- Verify email address is correct
- Check SMTP host is correct for your email provider

### SMTP Connection Failed

```
[ERROR] Could not connect to SMTP server
```

Solutions:
- Verify SMTP host (don't include port number)
- Check SMTP port is correct (usually 587)
- Try `UBOT_SMTP_INSECURE=true` if using non-standard certificates
- Ensure firewall allows outbound connections to SMTP port

### File Not Converting

```
[ERROR] could not convert file
```

Solutions:
- Check if file format is supported by Calibre
- Try a different file format
- Verify Calibre is installed in container:
  ```bash
  docker compose exec sendtokindle which ebook-convert
  ```

## Monitoring

### Real-time Logs

```bash
docker compose logs -f
```

### Check Resource Usage

```bash
docker stats sendtokindle-bot
```

### Container Health

```bash
docker inspect sendtokindle-bot --format='{{.State.Status}}'
```

## Updates and Maintenance

### Check for Updates

```bash
git status
git log --oneline -5
```

### Update to Latest Version

```bash
git pull origin master
docker compose build --no-cache
docker compose restart
```

### Backup Your Configuration

```bash
cp .env .env.backup
cp -r files files.backup
```

### Clean Up Old Images

```bash
docker image prune  # Remove unused images
docker system prune  # Remove unused containers, networks, volumes
```

## What Was Fixed

The latest rebuild includes these improvements:

✅ **File Extension Normalization** - Handles .PDF, .pdf, .Pdf equally
✅ **Email Subject Line** - Proper subject for Kindle delivery
✅ **SMTP Configuration** - Fixed port handling from environment variables
✅ **Better Error Messages** - Detailed logs for troubleshooting
✅ **Improved TLS Handling** - Correct SMTP authentication sequence
✅ **Path Handling** - More robust file system operations
✅ **Logging** - DEBUG, INFO, WARN, ERROR levels for monitoring

## Need Help?

Check the troubleshooting section in [README.md](README.md) for common issues.

Or review the full logs:
```bash
docker compose logs sendtokindle --tail 100
```
