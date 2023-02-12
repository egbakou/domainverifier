package domainverifier

import (
	"errors"
	"testing"
	"time"
)

func TestIsValidDomain(t *testing.T) {
	type args struct {
		domain string
	}
	testCases := []struct {
		name string
		args args
		want bool
	}{
		{"valid domain", args{"app-v1.fr.domain.live"}, true},
		{"invalid domain", args{"domain com"}, false},
		{"empty domain name", args{""}, false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidDomainName(tt.args.domain)
			if got != tt.want {
				t.Errorf("expected: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestIsSecure(t *testing.T) {
	type args struct {
		domain  string
		timeout time.Duration
	}
	testCases := []struct {
		name    string
		args    args
		want    bool
		wantErr error
	}{
		{"secure domain", args{"google.com", 5 * time.Second}, true, nil},
		{"insecure domain", args{"go.com", 5 * time.Second}, false, nil},
		{"invalid domain", args{"domain com", 5 * time.Second}, false,
			errors.New("invalid domain name")},
		{"unreachable domain", args{"unreachabledomain.com", 5 * time.Second}, false,
			errors.New("unreachable domain")},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			secure, err := IsSecure(tt.args.domain, tt.args.timeout)
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if secure != tt.want {
				t.Errorf("expected: %v, got: %v", tt.want, secure)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	testCases := []struct {
		args string
		want string
	}{
		{"   abc 123 ", "abc123"},
		{"abc@123", "abc123"},
		{"abc#123", "abc123"},
		{"Abc 123", "abc123"},
		{"Abc_123", "abc123"},
		{"Abc-123", "abc123"},
	}

	for _, tt := range testCases {
		got := sanitizeString(tt.args)
		if got != tt.want {
			t.Errorf("For input %q, expected %q, but got %q", tt.args, tt.want, got)
		}
	}
}
