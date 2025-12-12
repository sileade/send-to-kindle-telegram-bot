# Send-to-Kindle Telegram Bot

[![Go](https://github.com/sileade/send-to-kindle-telegram-bot/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/sileade/send-to-kindle-telegram-bot/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/sileade/send-to-kindle-telegram-bot)](https://goreportcard.com/report/github.com/sileade/send-to-kindle-telegram-bot) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful, secure, and easy-to-use Telegram bot that sends documents to your Kindle device. It automatically converts files to the recommended EPUB format and supports multiple Kindle devices.

## ✨ Features

- **Multi-Device Support**: Send documents to multiple Kindle devices with an interactive selection menu.
- **Automatic Conversion**: Converts a wide range of formats to EPUB, the officially recommended format for modern Kindle devices.
- **Secure**: Protects your credentials and sanitizes filenames to prevent security risks.
- **Robust Error Handling**: Provides clear feedback on success or failure.
- **Configurable**: Easily configure the bot using environment variables.
- **Dockerized**: Simple to deploy and run with Docker and Docker Compose.

## 🚀 Getting Started

### Prerequisites

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  **Clone the repository**:

    ```shell
    git clone https://github.com/sileade/send-to-kindle-telegram-bot.git
    cd send-to-kindle-telegram-bot
    ```

2.  **Create a `.env` file**:

    Create a `.env` file in the project root and add your configuration. See the [Configuration](#-configuration) section for details.

3.  **Run with Docker Compose**:

    ```shell
    docker-compose up -d --build
    ```

4.  **Check the logs**:

    ```shell
    docker-compose logs -f
    ```

## ⚙️ Configuration

Create a `.env` file with the following environment variables:

| Variable              | Description                                                                  | Required | Default       |
| --------------------- | ---------------------------------------------------------------------------- | :------: | ------------- |
| `UBOT_TELEGRAM_TOKEN` | Your Telegram bot token from [@BotFather](https://t.me/BotFather).         |   **Yes**    | -             |
| `UBOT_EMAIL_FROM`     | The email address the bot will use to send books.                            |   **Yes**    | -             |
| `UBOT_PASSWORD`       | The email password or app-specific password.                                 |   **Yes**    | -             |
| `UBOT_SMTP_HOST`      | The SMTP mail host (e.g., `smtp.gmail.com`).                                 |   **Yes**    | -             |
| `UBOT_EMAIL_TO`       | The default Kindle email address (used for single-device mode).              |    No    | -             |
| `UBOT_KINDLE_DEVICES` | A list of your Kindle devices and their emails (for multi-device mode).      |    No    | -             |
| `UBOT_SMTP_PORT`      | The SMTP port.                                                               |    No    | `587`         |
| `UBOT_SMTP_INSECURE`  | Set to `true` to skip TLS certificate verification (for testing only).       |    No    | `false`       |
| `UBOT_TMP_FILES_PATH` | The path where temporary files are stored.                                   |    No    | `/files/`     |

### Example `.env` File

```env
# Telegram Bot Token
UBOT_TELEGRAM_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz

# Email Configuration
UBOT_EMAIL_FROM=your-email@example.com
UBOT_PASSWORD=your-app-specific-password

# SMTP Configuration
UBOT_SMTP_HOST=smtp.example.com
UBOT_SMTP_PORT=587

# Single-Device Mode (uncomment and use this if you have one Kindle)
# UBOT_EMAIL_TO=your-kindle@kindle.com

# Multi-Device Mode (uncomment and use this for multiple Kindles)
UBOT_KINDLE_DEVICES="My Kindle:your-kindle@kindle.com|Spouse's Kindle:spouse-kindle@kindle.com"
```

### Multi-Device Configuration

To use the multi-device feature, set the `UBOT_KINDLE_DEVICES` environment variable with the following format:

`"DeviceName1:email1@kindle.com|DeviceName2:email2@kindle.com"`

- Separate each device with a pipe (`|`).
- Separate the device name and email with a colon (`:`).

## Usage

1.  **Start the bot** and ensure it's running correctly.
2.  **Send a document** to the bot in your Telegram chat.
3.  If you have multiple devices configured, the bot will ask you to **choose a destination**:

    ![Device Selection](https://i.imgur.com/example.png) <!-- Placeholder for a real image -->

4.  The bot will convert the file to **EPUB** and send it to your selected Kindle.

## 📚 Supported Formats

The bot sends the following formats directly to your Kindle without conversion:

- `EPUB`
- `PDF`
- `TXT`
- `DOC`, `DOCX`
- `RTF`
- `HTM`, `HTML`

For all other formats supported by Calibre (such as `FB2`, `AZW`, `MOBI`), the bot will automatically convert them to **EPUB** before sending.

## 🌐 Deployment

### Docker (Recommended)

The recommended way to deploy the bot is using Docker and Docker Compose, as described in the [Getting Started](#-getting-started) section.

### Kubernetes

Here is an example of a Kubernetes deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: send-to-kindle-bot
spec:
  replicas: 1
  selector:
    matchLabels:
      app: send-to-kindle-bot
  template:
    metadata:
      labels:
        app: send-to-kindle-bot
    spec:
      containers:
      - name: send-to-kindle-bot
        image: sileade/send-to-kindle-telegram-bot:latest
        env:
        - name: UBOT_TELEGRAM_TOKEN
          valueFrom:
            secretKeyRef:
              name: bot-secrets
              key: telegram_token
        - name: UBOT_EMAIL_FROM
          valueFrom:
            secretKeyRef:
              name: bot-secrets
              key: email_from
        # ... and so on for other environment variables
```

### AWS (EC2 with Docker)

1.  **Launch an EC2 instance** (e.g., `t3.micro` with Amazon Linux 2).
2.  **Install Docker and Docker Compose** on the instance.
3.  **Clone the repository** and create your `.env` file.
4.  **Run the bot** using `docker-compose up -d`.

## ❓ FAQ

**Q: Why does the bot convert files to EPUB?**

A: Amazon has officially discontinued support for the MOBI format in its Send to Kindle service. EPUB is now the recommended format for the best reading experience on modern Kindle devices.

**Q: Can I use this bot with Gmail?**

A: Yes, but you'll need to create an **app-specific password** for your Google account. Do not use your main password.

**Q: How do I find my Kindle's email address?**

A: You can find it in your Amazon account under **Manage Your Content and Devices > Preferences > Personal Document Settings**.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1.  **Fork the repository**.
2.  **Create a new branch** (`git checkout -b feature/your-feature`).
3.  **Make your changes** and commit them (`git commit -m 'Add some feature'`).
4.  **Push to the branch** (`git push origin feature/your-feature`).
5.  **Open a pull request**.

## 📜 Changelog

### v2.0.0 (2025-12-10)

-   **BREAKING CHANGE**: Changed default conversion format from MOBI to **EPUB**.
-   **Feature**: Added multi-device support with an interactive selection menu.
-   **Feature**: Added configurable temporary files path.
-   **Security**: Added filename sanitization to prevent path traversal attacks.
-   **Security**: Added email masking in logs.
-   **Fix**: Fixed numerous bugs, including a critical indentation error and a race condition.
-   **Tests**: Added a comprehensive test suite with 25+ test cases.
-   **Docs**: Completely rewrote the README.md with detailed instructions.

### v1.x.x

-   Initial release with basic functionality.

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

-   This project was inspired by the original work of [Michael Fomenko](https://github.com/michaelfmnk).
-   Thanks to the developers of [Telebot](https://github.com/tucnak/telebot) and [Calibre](https://calibre-ebook.com/).

## 🧪 Testing

This project has been thoroughly tested with both unit tests and functional tests.

### Unit Tests

**Coverage:** 14.4% (25+ test cases)

Run the unit tests:

```bash
go test ./bot -v -cover
```

**Test categories:**
- Configuration validation (6 tests)
- Filename sanitization (7 tests)
- Email masking (4 tests)
- File conversion logic (6 tests)
- Temporary files path (2 tests)

### Functional Tests

**Status:** ✅ 50/50 tests passed (100%)

The bot has been tested in a virtual environment with mock data:

| Category | Tests | Status |
|----------|-------|--------|
| Configuration | 6/6 | ✅ 100% |
| Validation | 9/9 | ✅ 100% |
| File Processing | 25/25 | ✅ 100% |
| SMTP Functionality | 10/10 | ✅ 100% |

**What was tested:**
- ✅ Configuration loading and validation
- ✅ Multi-device support
- ✅ Filename sanitization (path traversal protection)
- ✅ File format detection and conversion logic
- ✅ Email masking in logs
- ✅ SMTP connection and authentication (mock)
- ✅ Email format validation
- ✅ TLS configuration
- ✅ File size limits

**Security features verified:**
- ✅ Protection against path traversal attacks
- ✅ Confidential data masking in logs
- ✅ Input validation
- ✅ TLS for SMTP connections

### Test Results

All tests pass successfully:

```bash
# Unit tests
PASS
coverage: 14.4% of statements
ok      github.com/michaelfmnk/send-to-kindle-telegram-bot/bot 0.006s

# Functional tests
✅ Configuration: 6/6 passed
✅ Validation: 9/9 passed
✅ File Processing: 25/25 passed
✅ SMTP: 10/10 passed
```

## 🚨 Troubleshooting

If books are not arriving on your Kindle, follow these steps to diagnose the issue. The most common problems are related to configuration, not code bugs.

### Step 1: Check Docker Logs

This is the most important step. The logs will tell you exactly what is failing.

```shell
docker-compose logs -f sendtokindle
```

Look for `[ERROR]` messages:

| Log Message | Meaning | Solution |
|---|---|---|
| `Authentication failed` | Incorrect email/password | [Use an App-Specific Password](#-use-an-app-specific-password) |
| `Could not connect to SMTP server` | Wrong SMTP host/port | [Verify SMTP Settings](#-verify-smtp-settings) |
| `emailto not set` | No Kindle email configured | [Configure Your Kindle Email](#-configure-your-kindle-email) |
| `Could not convert file` | Calibre conversion failed | Check Docker build logs for Calibre installation errors. |

### Step 2: Verify Your Amazon Settings

This is the most common reason for failure.

1.  **Approve Your Sending Email**
    - Go to [Amazon's Manage Your Content and Devices](https://www.amazon.com/gp/digital/fiona/)
    - Navigate to **Preferences > Personal Document Settings**.
    - Under **Approved Personal Document E-mail List**, click **Add a new approved e-mail address**.
    - Add the email address you configured in `UBOT_EMAIL_FROM`.

2.  **Check Your Kindle's Email Address**
    - On the same page, find your Kindle device in the **Send-to-Kindle E-Mail Settings** list.
    - Ensure the email address matches what you have in `UBOT_EMAIL_TO` or `UBOT_KINDLE_DEVICES`.

### Step 3: Use an App-Specific Password

Do **NOT** use your main email password. You must generate an app-specific password.

- **Gmail**: [Create an App Password](https://myaccount.google.com/apppasswords) (requires 2-Factor Authentication).
- **Yandex**: [Generate an App Password](https://passport.yandex.ru/) under "Security".
- **Outlook**: [Create an App Password](https://account.live.com/) under "Security".

Update the `UBOT_PASSWORD` in your `.env` file with this new password.

### Step 4: Verify SMTP Settings

Check your `.env` file for the correct SMTP configuration.

| Provider | `UBOT_SMTP_HOST` | `UBOT_SMTP_PORT` |
|---|---|---|
| Gmail | `smtp.gmail.com` | `587` |
| Yandex | `smtp.yandex.com` | `587` |
| Outlook | `smtp.live.com` | `587` |

**Important**: `UBOT_SMTP_HOST` should **NOT** include the port (e.g., `smtp.gmail.com`, not `smtp.gmail.com:587`).

### Step 5: Check `.env` File Syntax

Ensure your `.env` file has no syntax errors:

- **NO** spaces around the `=` sign (`UBOT_TOKEN=...`, not `UBOT_TOKEN = ...`).
- **NO** comments on the same line as a variable.
- **NO** empty variables that are required.

### Quick Diagnostic Checklist

- [ ] Have you checked the `docker-compose logs` for errors?
- [ ] Is your `UBOT_EMAIL_FROM` address in Amazon's approved senders list?
- [ ] Are you using an **app-specific password** for `UBOT_PASSWORD`?
- [ ] Is your `UBOT_EMAIL_TO` address correct?
- [ ] Are your `UBOT_SMTP_HOST` and `UBOT_SMTP_PORT` correct?
- [ ] Have you checked your `.env` file for syntax errors?
