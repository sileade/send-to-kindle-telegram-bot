package bot

// SetTmpFilesPath sets the temporary files path for the bot
func (b *SendToKindleBot) SetTmpFilesPath(path string) {
	if path != "" {
		b.tmpFilesPath = path
	}
}

// GetTmpFilesPath returns the temporary files path
func (b *SendToKindleBot) GetTmpFilesPath() string {
	if b.tmpFilesPath == "" {
		return defaultTmpFilesPath
	}
	return b.tmpFilesPath
}
