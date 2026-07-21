package cmd

import (
	"path/filepath"
	"testing"
)

func TestResolveDownloadDestPath(t *testing.T) {
	tests := []struct {
		name      string
		dest      string
		destIsDir bool
		domain    string
		certID    string
		want      string
	}{
		{
			name:      "current dir default builds domain-id zip",
			dest:      ".",
			destIsDir: true,
			domain:    "www.example.com",
			certID:    "abc123",
			want:      filepath.Join(".", "www.example.com-abc123.zip"),
		},
		{
			name:      "directory builds domain-id zip",
			dest:      "/tmp/certs",
			destIsDir: true,
			domain:    "www.example.com",
			certID:    "abc123",
			want:      filepath.Join("/tmp/certs", "www.example.com-abc123.zip"),
		},
		{
			name:      "directory with trailing slash",
			dest:      "certs/",
			destIsDir: true,
			domain:    "a.io",
			certID:    "x",
			want:      filepath.Join("certs", "a.io-x.zip"),
		},
		{
			name:      "directory sanitizes wildcard domain",
			dest:      "out",
			destIsDir: true,
			domain:    "_.example.com",
			certID:    "id1",
			want:      filepath.Join("out", "_.example.com-id1.zip"),
		},
		{
			name:      "directory falls back to id-only when domain unknown",
			dest:      "out",
			destIsDir: true,
			domain:    "",
			certID:    "id1",
			want:      filepath.Join("out", "id1.zip"),
		},
		{
			name:      "file with .zip is used verbatim",
			dest:      "cert.zip",
			destIsDir: false,
			want:      "cert.zip",
		},
		{
			name:      "file with .ZIP (uppercase) is used verbatim",
			dest:      "cert.ZIP",
			destIsDir: false,
			want:      "cert.ZIP",
		},
		{
			name:      "non-zip extension gets .zip appended",
			dest:      "cert.pem",
			destIsDir: false,
			want:      "cert.pem.zip",
		},
		{
			name:      "bare name gets .zip appended",
			dest:      "cert",
			destIsDir: false,
			want:      "cert.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveDownloadDestPath(tt.dest, tt.destIsDir, tt.domain, tt.certID)
			if got != tt.want {
				t.Errorf("resolveDownloadDestPath(%q, %v, %q, %q) = %q, want %q",
					tt.dest, tt.destIsDir, tt.domain, tt.certID, got, tt.want)
			}
		})
	}
}

func TestCanonicalDomain(t *testing.T) {
	tests := []struct {
		name string
		sans []string
		want string
	}{
		{name: "first san", sans: []string{"www.example.com", "www-alt.example.com"}, want: "www.example.com"},
		{name: "wildcard sanitized", sans: []string{"*.example.com"}, want: "_.example.com"},
		{name: "empty", sans: []string{}, want: ""},
		{name: "nil", sans: nil, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canonicalDomain(tt.sans); got != tt.want {
				t.Errorf("canonicalDomain(%v) = %q, want %q", tt.sans, got, tt.want)
			}
		})
	}
}
