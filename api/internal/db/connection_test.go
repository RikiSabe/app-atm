package db

import (
	"net/url"
	"testing"
)

func TestSupabaseDSNRequiresTLS(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		wantSSLMode string
	}{
		{
			name:        "adds TLS when no mode is provided",
			input:       "postgresql://user:secret@db.example.supabase.co:5432/postgres",
			wantSSLMode: "require",
		},
		{
			name:        "preserves TLS verification mode",
			input:       "postgres://user:secret@db.example.supabase.co:5432/postgres?sslmode=verify-full",
			wantSSLMode: "verify-full",
		},
		{
			name:    "rejects disabled TLS",
			input:   "postgres://user:secret@db.example.supabase.co:5432/postgres?sslmode=disable",
			wantErr: true,
		},
		{
			name:    "rejects a missing URL",
			wantErr: true,
		},
		{
			name:    "rejects a non-PostgreSQL URL",
			input:   "https://example.supabase.co",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := supabaseDSN(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("supabaseDSN() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}

			parsed, err := url.Parse(got)
			if err != nil {
				t.Fatalf("supabaseDSN() returned an invalid URL: %v", err)
			}
			if gotSSLMode := parsed.Query().Get("sslmode"); gotSSLMode != test.wantSSLMode {
				t.Errorf("sslmode = %q, want %q", gotSSLMode, test.wantSSLMode)
			}
		})
	}
}
