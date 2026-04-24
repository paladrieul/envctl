package store

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestExportDotenv(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("prod", "DB_HOST", "localhost")
	_ = s.Set("prod", "APP_NAME", "my app")

	r, w, _ := os.Pipe()
	err := s.Export("prod", FormatDotenv, w)
	w.Close()
	if err != nil {
		t.Fatalf("Export dotenv: %v", err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME=") {
		t.Errorf("expected APP_NAME in output, got:\n%s", out)
	}
}

func TestExportJSON(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("staging", "PORT", "8080")
	_ = s.Set("staging", "DEBUG", "true")

	r, w, _ := os.Pipe()
	err := s.Export("staging", FormatJSON, w)
	w.Close()
	if err != nil {
		t.Fatalf("Export JSON: %v", err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal JSON output: %v", err)
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result["PORT"])
	}
	if result["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", result["DEBUG"])
	}
}

func TestExportShell(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("dev", "HOME", "/home/user")

	r, w, _ := os.Pipe()
	err := s.Export("dev", FormatExport, w)
	w.Close()
	if err != nil {
		t.Fatalf("Export shell: %v", err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "export HOME=/home/user") {
		t.Errorf("expected 'export HOME=/home/user' in output, got:\n%s", out)
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("dev", "KEY", "val")
	err := s.Export("dev", ExportFormat("xml"), nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExportNonExistentTarget(t *testing.T) {
	s := tempStore(t)
	r, w, _ := os.Pipe()
	err := s.Export("ghost", FormatDotenv, w)
	w.Close()
	r.Close()
	if err != nil {
		t.Fatalf("exporting non-existent target should return empty, got: %v", err)
	}
}
