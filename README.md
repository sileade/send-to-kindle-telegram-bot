# SendToKindle [![Go](https://github.com/michaelfmnk/SendToKindleTelegramBot/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/michaelfmnk/SendToKindleTelegramBot/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/michaelfmnk/send-to-kindle-telegram-bot)](https://goreportcard.com/report/github.com/michaelfmnk/send-to-kindle-telegram-bot)

## Building the Docker Container

To build the Docker container, execute the following command:

```shell
docker build -t sendtokindle .
```

Or using docker-compose:

```shell
docker compose build
```

## Configuring Environment Variables

Before starting the bot, you need to configure the following environment variables:

| Variable              | Description                                                      | Required | Default |
|-----------------------|------------------------------------------------------------------|----------|----------|
| UBOT_TELEGRAM_TOKEN   | Telegram bot token                                               | Yes      | -       |
| UBOT_EMAIL_FROM       | Email address that the bot will use to send books                | Yes      | -       |
| UBOT_PASSWORD         | Email password or app-specific password                          | Yes      | -       |
| UBOT_EMAIL_TO         | Kindle email address to which books will be sent                 | Yes      | -       |
| UBOT_SMTP_HOST        | SMTP mail host (without port)                                    | Yes      | -       |
| UBOT_SMTP_PORT        | SMTP port                                                        | No       | 587     |
| UBOT_SMTP_INSECURE    | Skip TLS certificate verification (use "true" or "1" to enable) | No       | false   |

### Example .env file

```bash
UBOT_TELEGRAM_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
UBOT_EMAIL_FROM=your-email@example.com
UBOT_PASSWORD=your-password
UBOT_EMAIL_TO=your-kindle@kindle.com
UBOT_SMTP_HOST=smtp.example.com
# Optional: set to true for self-signed certificates
UBOT_SMTP_INSECURE=false
# Optional: specify SMTP port if different from 587
UBOT_SMTP_PORT=587
```

## Usage

After starting the bot and configuring the necessary environment variables, you can use it by sending a message to the bot containing the document you want to convert. The bot will then send the converted document to your Kindle email address.

### Using Docker Compose

1. Create a `.env` file with your configuration (see example above)
2. Create a `docker-compose.yml` file:

```yaml
services:
  sendtokindle:
    build: .
    container_name: sendtokindle-bot
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - ./files:/files
```

3. Run the bot:

```bash
docker compose up -d
```

4. Check logs:

```bash
docker compose logs -f
```

## Supported Formats

The bot will automatically convert the following formats to MOBI:
- EPUB
- FB2
- AZW
- And many others supported by Calibre

These formats are sent directly without conversion:
- DOC, DOCX
- RTF
- HTM, HTML
- TXT
- MOBI
- PDF

## Troubleshooting

### Books are not being delivered to Kindle

#### 1. Check Bot Logs
Always start by reviewing the bot logs for errors:
```bash
docker compose logs -f sendtokindle
```

Look for `[ERROR]` messages that indicate what went wrong.

#### 2. Verify Email Configuration
- **Email Address**: Make sure `UBOT_EMAIL_FROM` is a valid email address
- **Password**: Use an app-specific password, not your main email password
  - Gmail: [Create app password](https://support.google.com/accounts/answer/185833)
  - Yandex: [Generate app password](https://yandex.com/support/id/authorization/app-passwords.html)
  - Outlook: [Create app password](https://support.microsoft.com/en-us/account-billing/using-app-passwords-with-your-microsoft-account)

#### 3. Verify Kindle Email Settings
- **Kindle Email**: Must be registered in [Amazon Manage Your Content and Devices](https://www.amazon.com/gp/digital/fiona/)
- **Email Whitelist**: Add your `UBOT_EMAIL_FROM` to the [Approved Personal Document E-mail List](https://www.amazon.com/gp/digital/fiona/)
- **Receive Documents**: Enable "Receive e-mail-based Kindle documents"

#### 4. SMTP Configuration Issues
- **SMTP Host**: Use the full host without port (e.g., `smtp.gmail.com`, not `smtp.gmail.com:587`)
- **SMTP Port**: Common ports are:
  - 587 (TLS, recommended)
  - 465 (SSL)
  - 25 (rarely used now)
- **TLS Certificate**: If using a self-signed certificate, set `UBOT_SMTP_INSECURE=true`

#### 5. File Format Issues
- **Case Sensitivity**: File extensions are now normalized to lowercase (e.g., `.PDF` will work)
- **Calibre Conversion**: If conversion fails, check if calibre is properly installed (it's included in Docker)
- **File Size**: Amazon Kindle has limits on document size (typically 25-50MB)

#### 6. Check Docker Container
```bash
# Verify container is running
docker ps

# Check logs for startup errors
docker compose logs sendtokindle

# Rebuild if needed
docker compose down
docker compose build --no-cache
docker compose up -d
```

### Common Error Messages

**"could not connect to SMTP server"**
- Verify SMTP host is correct
- Check SMTP port
- Ensure firewall allows outbound connections

**"authentication failed"**
- Check email address and password/app-password
- Verify email account allows SMTP connections
- Try `UBOT_SMTP_INSECURE=true` if using non-standard certificates

**"could not convert file"**
- Ensure file format is supported by Calibre
- Check Docker logs for specific conversion errors
- Try manually converting with: `docker exec sendtokindle-bot ebook-convert <input> <output.mobi>`

## Links

For more information, you can check out the following links:

- [Dev.to post](https://dev.to/michaelfmnk/developing-send-to-kindle-telegram-bot-120c)
- [Amazon Kindle Document Support](https://www.amazon.com/gp/help/customer/display.html?nodeId=GKMQC26VQQMM8XKW)
- If you have any questions or suggestions, feel free to message Michael at [michael@fomenko.dev](mailto:michael@fomenko.dev).
