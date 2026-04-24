package store

import (
	"strings"
	"testing"
)

func TestImportDotenv(t *testing.T) {
	s := tempStore(t)
	input := `# comment
DB_HOST=localhost
DB_PORT="5432"
SECRET=abc123
`
	res, err := Import(s, "prod", strings.NewReader(input), FormatDotenv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 3 {
		t.Errorf("expected 3 imported, got %d", res.Imported)
	}

	env, _ := s.Load("prod")
	if env["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q, want %q", env["DB_HOST"], "localhost")
	}
	if env["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: got %q, want %q", env["DB_PORT"], "5432")
	}
}

func TestImportJSON(t *testing.T) {
	s := tempStore(t)
	input := `{"APP_ENV":"staging","LOG_LEVEL":"debug"}`
	res, err := Import(s, "staging", strings.NewReader(input), FormatJSON, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 2 {
		t.Errorf("expected 2 imported, got %d", res.Imported)
	}

	env, _ := s.Load("staging")
	if env["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV: got %q", env["APP_ENV"])
	}
}

func TestImportSkipsExistingWithoutOverwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("dev", map[string]string{"KEY": "original"})

	input := "KEY=new\nOTHER=value\n"
	res, err := Import(s, "dev", strings.NewReader(input), FormatDotenv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 || res.Imported != 1 {
		t.Errorf("expected 1 skipped, 1 imported; got skipped=%d imported=%d", res.Skipped, res.Imported)
	}

	env, _ := s.Load("dev")
	if env["KEY"] != "original" {
		t.Errorf("KEY should not have been overwritten, got %q", env["KEY"])
	}
}

func TestImportOverwritesExisting(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("dev", map[string]string{"KEY": "original"})

	input := "KEY=updated\n"
	res, err := Import(s, "dev", strings.NewReader(input), FormatDotenv, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 1 {
		t.Errorf("expected 1 imported, got %d", res.Imported)
	}

	env, _ := s.Load("dev")
	if env["KEY"] != "updated" {
		t.Errorf("KEY should be updated, got %q", env["KEY"])
	}
}

func TestImportUnsupportedFormat(t *testing.T) {
	s := tempStore(t)
	_, err := Import(s, "prod", strings.NewReader(""), "yaml", false)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
