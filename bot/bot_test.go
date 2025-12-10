package bot

import (
	"testing"
)

func TestSendToKindleBot_verifyConfig(t *testing.T) {
	type fields struct {
		Token     string
		EmailFrom string
		EmailTo   string
		SMTPHost  string
		SMTPPort  string
		Password  string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		errorMessage string
	}{
		{
			name: "should pass validation with all required fields",
			fields: fields{
				Token:     "hello",
				EmailFrom: "from@example.com",
				EmailTo:   "to@kindle.com",
				SMTPHost:  "smtp.example.com",
				SMTPPort:  "587",
				Password:  "pass",
			},
			wantErr: false,
		},
		{
			name: "should not pass validation if Token empty",
			fields: fields{
				Token:     "",
				EmailFrom: "from@example.com",
				EmailTo:   "to@kindle.com",
				SMTPHost:  "smtp.example.com",
				SMTPPort:  "587",
				Password:  "pass",
			},
			wantErr:      true,
			errorMessage: "token for telegram bot not set",
		},
		{
			name: "should not pass validation if Password empty",
			fields: fields{
				Token:     "hello",
				EmailFrom: "from@example.com",
				EmailTo:   "to@kindle.com",
				SMTPHost:  "smtp.example.com",
				SMTPPort:  "587",
				Password:  "",
			},
			wantErr:      true,
			errorMessage: "password for email not set",
		},
		{
			name: "should not pass validation if EmailFrom empty",
			fields: fields{
				Token:     "hello",
				EmailFrom: "",
				EmailTo:   "to@kindle.com",
				SMTPHost:  "smtp.example.com",
				SMTPPort:  "587",
				Password:  "pass",
			},
			wantErr:      true,
			errorMessage: "emailfrom not set",
		},
		{
			name: "should not pass validation if SMTPHost empty",
			fields: fields{
				Token:     "hello",
				EmailFrom: "from@example.com",
				EmailTo:   "to@kindle.com",
				SMTPHost:  "",
				SMTPPort:  "587",
				Password:  "pass",
			},
			wantErr:      true,
			errorMessage: "smtp host not set",
		},
		{
			name: "should use default SMTP port if not provided",
			fields: fields{
				Token:     "hello",
				EmailFrom: "from@example.com",
				EmailTo:   "to@kindle.com",
				SMTPHost:  "smtp.example.com",
				SMTPPort:  "",
				Password:  "pass",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &SendToKindleBot{
				Token:     tt.fields.Token,
				EmailFrom: tt.fields.EmailFrom,
				EmailTo:   tt.fields.EmailTo,
				SMTPHost:  tt.fields.SMTPHost,
				SMTPPort:  tt.fields.SMTPPort,
				Password:  tt.fields.Password,
			}
			err := b.verifyConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && err.Error() != tt.errorMessage {
				t.Errorf("verifyConfig() error '%v' not matches '%s'", err, tt.errorMessage)
			}

			// Check default port was set
			if !tt.wantErr && tt.fields.SMTPPort == "" && b.SMTPPort != defaultSMTPPort {
				t.Errorf("verifyConfig() expected default SMTP port %s, got %s", defaultSMTPPort, b.SMTPPort)
			}
		})
	}
}

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		wantErr   bool
		wantName  string
		wantEqual bool
	}{
		{
			name:      "valid filename",
			fileName:  "document.pdf",
			wantErr:   false,
			wantName:  "document.pdf",
			wantEqual: true,
		},
		{
			name:      "filename with spaces",
			fileName:  "my document.pdf",
			wantErr:   false,
			wantName:  "my document.pdf",
			wantEqual: true,
		},
		{
			name:      "filename with path traversal attempt",
			fileName:  "../../../etc/passwd",
			wantErr:   false,
			wantName:  "passwd",
			wantEqual: true,
		},
		{
			name:     "empty filename",
			fileName: "",
			wantErr:  true,
		},
		{
			name:     "filename with only dots",
			fileName: "...",
			wantErr:  false,
			wantName: ".",
		},
		{
			name:      "filename with invalid characters",
			fileName:  "file<>name.pdf",
			wantErr:   false,
			wantName:  "filename.pdf",
			wantEqual: true,
		},
		{
			name:     "filename too long",
			fileName: string(make([]byte, 300)),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeFileName(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("sanitizeFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.wantEqual && got != tt.wantName {
				t.Errorf("sanitizeFileName() got = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{
			name:  "valid email",
			email: "user@example.com",
			want:  "***@example.com",
		},
		{
			name:  "kindle email",
			email: "user123@kindle.com",
			want:  "***@kindle.com",
		},
		{
			name:  "invalid email without @",
			email: "notanemail",
			want:  "***@***",
		},
		{
			name:  "invalid email with multiple @",
			email: "user@domain@example.com",
			want:  "***@***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maskEmail(tt.email)
			if got != tt.want {
				t.Errorf("maskEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeedToConvert(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		want      bool
	}{
		{
			name:      "epub format (no conversion needed)",
			extension: "epub",
			want:      false,
		},
		{
			name:      "pdf format (no conversion needed)",
			extension: "pdf",
			want:      false,
		},
		{
			name:      "txt format (no conversion needed)",
			extension: "txt",
			want:      false,
		},
		{
			name:      "fb2 format (conversion needed)",
			extension: "fb2",
			want:      true,
		},
		{
			name:      "azw format (conversion needed)",
			extension: "azw",
			want:      true,
		},
		{
			name:      "unknown format (conversion needed)",
			extension: "xyz",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needToConvert(tt.extension)
			if got != tt.want {
				t.Errorf("needToConvert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendToKindleBot_GetTmpFilesPath(t *testing.T) {
	tests := []struct {
		name     string
		bot      *SendToKindleBot
		wantPath string
	}{
		{
			name: "custom path set",
			bot: &SendToKindleBot{
				tmpFilesPath: "/custom/path",
			},
			wantPath: "/custom/path",
		},
		{
			name: "default path when empty",
			bot: &SendToKindleBot{
				tmpFilesPath: "",
			},
			wantPath: defaultTmpFilesPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bot.GetTmpFilesPath()
			if got != tt.wantPath {
				t.Errorf("GetTmpFilesPath() got = %v, want %v", got, tt.wantPath)
			}
		})
	}
}
