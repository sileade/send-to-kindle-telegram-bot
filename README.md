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
|-----------------------|------------------------------------------------------------------|----------|--------|
| UBOT_TELEGRAM_TOKEN   | Telegram bot token                                               | Yes      | -      |
| UBOT_EMAIL_FROM       | Email address that the bot will use to send books                | Yes      | -      |
| UBOT_PASSWORD         | Email password or app-specific password                          | Yes      | -      |
| UBOT_EMAIL_TO         | Kindle email address to which books will be sent                 | Yes      | -      |
| UBOT_SMTP_HOST        | SMTP mail host (without port)                                    | Yes      | -      |
| UBOT_SMTP_PORT        | SMTP port                                                        | No       | 587    |
| UBOT_SMTP_INSECURE    | Skip TLS certificate verification (use "true" or "1" to enable) | No       | false  |

### Example .env file

```bash
UBOT_TELEGRAM_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
UBOT_EMAIL_FROM=your-email@example.com
UBOT_PASSWORD=your-password
UBOT_EMAIL_TO=your-kindle@kindle.com
UBOT_SMTP_HOST=smtp.example.com
# Optional: set to true for self-signed certificates
UBOT_SMTP_INSECURE=true
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

## Links

For more information, you can check out the following links:

- [Dev.to post](https://dev.to/michaelfmnk/developing-send-to-kindle-telegram-bot-120c)
- If you have any questions or suggestions, feel free to message Michael at [michael@fomenko.dev](mailto:michael@fomenko.dev).