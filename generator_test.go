package domainverifier

import (
	"errors"
	"github.com/egbakou/domainverifier/config"
	"strings"
	"testing"
)

func TestGenerateHtmlMetaFromConfig(t *testing.T) {
	type args struct {
		config          *config.HmlMetaTagGenerator
		useInternalCode bool
	}
	testCases := []struct {
		name      string
		args      args
		want      *HtmlMetaInstruction
		wantError error
	}{
		{
			name: "Successful generation with internal code",
			args: args{
				config: &config.HmlMetaTagGenerator{
					TagName: "example-tag",
				},
				useInternalCode: true,
			},
			want: &HtmlMetaInstruction{
				Code: `<meta name="example-tag" content=`,
			},
			wantError: nil,
		},
		{
			name: "Successful generation with external code",
			args: args{config: &config.HmlMetaTagGenerator{
				TagName: "example-tag",
				Code:    "external-code",
			},
				useInternalCode: false,
			},
			want: &HtmlMetaInstruction{
				Code: `<meta name="example-tag" content="external-code" />`,
			},
			wantError: nil,
		},
		{
			name: "Validation error",
			args: args{
				config:          nil,
				useInternalCode: true,
			},
			want:      nil,
			wantError: config.InvalidConfigError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateHtmlMetaFromConfig(tt.args.config, tt.args.useInternalCode)
			if err == nil && tt.wantError != nil {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}

			if tt.wantError == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if got != nil && !strings.Contains(got.Code, tt.want.Code) {
				t.Errorf("expected: %v, got: %v", tt.want.Code, got.Code)
			}
		})
	}
}

func TestGenerateHtmlMeta(t *testing.T) {
	type args struct {
		appName  string
		sanitize bool
	}
	testCases := []struct {
		name    string
		args    args
		want    *HtmlMetaInstruction
		wantErr error
	}{
		{
			name: "valid app name and sanitize is false",
			args: args{
				appName:  "my super app",
				sanitize: false,
			},
			want: &HtmlMetaInstruction{
				Code: `<meta name="my super app" content=`,
			},
			wantErr: nil,
		},
		{
			name: "valid app name and sanitize is true",
			args: args{
				appName:  "my super app",
				sanitize: true,
			},
			want: &HtmlMetaInstruction{
				Code: `<meta name="mysuperapp" content=`,
			},
			wantErr: nil,
		},
		{
			name: "invalid app name",
			args: args{
				appName:  "",
				sanitize: false,
			},
			want:    nil,
			wantErr: InvalidAppNameError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateHtmlMeta(tt.args.appName, tt.args.sanitize)
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if got != nil && !strings.Contains(got.Code, tt.want.Code) {
				t.Errorf("expected: %v, got: %v", tt.want.Code, got.Code)
			}
		})
	}
}

func TestGenerateJsonFromConfig(t *testing.T) {
	type args struct {
		config          *config.JsonGenerator
		useInternalCode bool
	}
	tests := []struct {
		name    string
		args    args
		want    *FileInstruction
		wantErr bool
	}{
		{
			name: "Successful generation with external code",
			args: args{
				config: &config.JsonGenerator{
					FileName:  "example.json",
					Attribute: "code",
					Code:      "external-code",
				},
				useInternalCode: false,
			},
			want: &FileInstruction{
				FileName:    "example.json",
				FileContent: `{"code": "external-code"}`,
			},
			wantErr: false,
		},
		{
			name: "Successful generation with internal code",
			args: args{
				config: &config.JsonGenerator{
					FileName:  "example",
					Attribute: "code",
				},
				useInternalCode: true,
			},
			want: &FileInstruction{
				FileName:    "example",
				FileContent: `{"code": "`,
			},
			wantErr: false,
		},
		{
			name: "Validation error",
			args: args{
				config:          nil,
				useInternalCode: true,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJsonFromConfig(tt.args.config, tt.args.useInternalCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJsonFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !strings.Contains(got.FileContent, tt.want.FileContent) {
				t.Errorf("expected: %v, got: %v", tt.want.FileContent, got.FileContent)
			}
		})
	}
}

func TestGenerateJson(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    *FileInstruction
		wantErr bool
	}{
		{
			name: "valid app name and sanitize is false",
			args: "my super app",
			want: &FileInstruction{
				FileName:    "my super app.json",
				FileContent: `{"mysuperapp_site_verification": "`,
			},
			wantErr: false,
		},
		{
			name: "valid app name and sanitize is true",
			args: "my super app",
			want: &FileInstruction{
				FileName:    "mysuperapp.json",
				FileContent: `{"mysuperapp_site_verification": "`,
			},
			wantErr: false,
		},
		{
			name:    "invalid app name",
			args:    "",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJson(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !strings.Contains(got.FileContent, tt.want.FileContent) {
				t.Errorf("expected: %v, got: %v", tt.want.FileContent, got.FileContent)
			}
		})
	}
}

func TestGenerateXmlFromConfig(t *testing.T) {
	type args struct {
		config          *config.XmlGenerator
		useInternalCode bool
	}
	testCases := []struct {
		name    string
		args    args
		want    *FileInstruction
		wantErr error
	}{
		{
			name: "Successful generation with external code",
			args: args{
				config: &config.XmlGenerator{
					FileName: "example.xml",
					RootName: "example-root",
					Code:     "internal-code",
				},
				useInternalCode: false,
			},
			want: &FileInstruction{
				FileName:    "example.xml",
				FileContent: `<example-root><code>internal-code</code></example-root>`,
			},
			wantErr: nil,
		},
		{
			name: "Successful generation with internal code",
			args: args{
				config: &config.XmlGenerator{
					FileName: "example",
					RootName: "example-root",
				},
				useInternalCode: true,
			},
			want: &FileInstruction{
				FileName:    "example",
				FileContent: `<example-root><code>`,
			},
			wantErr: nil,
		},
		{
			name: "Validation error",
			args: args{
				config:          nil,
				useInternalCode: true,
			},
			want:    nil,
			wantErr: errors.New("config cannot be nil"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateXmlFromConfig(tt.args.config, tt.args.useInternalCode)
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if got != nil && !strings.Contains(got.FileContent, tt.want.FileContent) {
				t.Errorf("expected: %v, got: %v", tt.want.FileContent, got.FileContent)
			}
		})
	}
}

func TestGenerateXml(t *testing.T) {
	type args struct {
		appName         string
		sanitizeAppName bool
	}
	tests := []struct {
		name    string
		args    args
		want    *FileInstruction
		wantErr error
	}{
		{
			name: "valid app name and sanitize is false",
			args: args{
				appName:         "my super app",
				sanitizeAppName: false,
			},
			want: &FileInstruction{
				FileName:    "my super app.xml",
				FileContent: `<verification><code>`,
			},
			wantErr: nil,
		},
		{
			name: "valid app name and sanitize is true",
			args: args{
				appName:         "my super app",
				sanitizeAppName: true,
			},
			want: &FileInstruction{
				FileName:    "mysuperapp.xml",
				FileContent: `<verification><code>`,
			},
			wantErr: nil,
		},
		{
			name: "invalid app name",
			args: args{
				appName: "",
			},
			want:    nil,
			wantErr: InvalidAppNameError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateXml(tt.args.appName, tt.args.sanitizeAppName)
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if got != nil && !strings.Contains(got.FileContent, tt.want.FileContent) {
				t.Errorf("expected: %v, got: %v", tt.want.FileContent, got.FileContent)
			}
		})
	}
}
