# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- 🔧 **File Extension Normalization**: Extensions are now normalized to lowercase, so `.PDF`, `.pdf`, `.Pdf` are all handled consistently
- 📧 **Email Subject Line**: Emails now include meaningful subject lines with book names for better Kindle delivery
- 📝 **Comprehensive Logging**: Added DEBUG, INFO, WARN, ERROR level logging for easier troubleshooting
- 🐳 **Docker Compose Support**: Added `docker-compose.yml` with proper configuration and logging
- 📚 **Environment Configuration Example**: Added `.env.example` with detailed documentation
- 🛠️ **Rebuild Guide**: New `REBUILD.md` with step-by-step deployment instructions
- 🆘 **Troubleshooting Guide**: Extended README with comprehensive troubleshooting section
- ✅ **Path Handling Improvements**: More robust file path handling using `filepath` package
- 🔐 **Enhanced SMTP TLS**: Improved error messages for SMTP/TLS connection issues

### Fixed
- 🐛 **Missing SMTP Port Configuration**: `UBOT_SMTP_PORT` environment variable is now properly used
- 🐛 **Directory Creation**: Automatic creation of `/files/` directory if it doesn't exist
- 🐛 **Email Subject**: Emails now have proper subject lines instead of empty subjects
- 🐛 **Error Messages**: Detailed and descriptive error messages for debugging
- 🐛 **Authentication Error Handling**: Clear error messages for SMTP authentication failures
- 🐛 **SMTP Host Validation**: Added validation for SMTP host configuration

### Changed
- 📄 **Logging Format**: Changed from generic log.Println to structured logging with [LEVEL] tags
- 🔄 **Error Wrapping**: Better error context using `fmt.Errorf` with `%w` for error wrapping
- 🎯 **User Feedback**: More descriptive messages sent to Telegram users
- 📖 **Documentation**: Significantly expanded README with troubleshooting and setup guides

### Technical Details

#### bot/bot.go
- Added `filepath.Ext()` for robust file extension extraction
- Added file extension normalization with `strings.ToLower()`
- Improved `sendFileViaEmail()` signature to include original filename
- Enhanced SMTP TLS error handling with detailed error messages
- Added debug logging at key points in the process
- Improved path handling using `filepath.Join()`
- Added `ensureDirectory()` function for directory creation
- Added `ErrNoSMTPHost` error type

#### main.go
- Added `UBOT_SMTP_PORT` from environment variables
- Improved error logging for startup failures

#### Configuration Files
- Updated `docker-compose.yml` with logging configuration
- Created `.env.example` with comprehensive documentation
- Added resource limits configuration (commented out for flexibility)

## How to Update

1. Pull the latest changes:
   ```bash
   git pull origin master
   ```

2. Rebuild the Docker container:
   ```bash
   docker compose build --no-cache
   docker compose up -d
   ```

3. Check logs to verify everything is working:
   ```bash
   docker compose logs -f
   ```

See [REBUILD.md](REBUILD.md) for detailed deployment instructions.

## Support

If you encounter any issues:

1. Check the [Troubleshooting Guide](README.md#troubleshooting) in README.md
2. Review the detailed logs: `docker compose logs sendtokindle`
3. Verify your `.env` configuration matches the requirements

## Commits

- ✅ fix: improve email sending, add file extension normalization, and fix SMTP configuration issues
- ✅ fix: add missing UBOT_SMTP_PORT environment variable configuration
- ✅ docs: update README with improvements and troubleshooting section
- ✅ feat: add docker-compose configuration with improvements
- ✅ docs: add environment configuration example
- ✅ docs: add container rebuild and deployment guide
