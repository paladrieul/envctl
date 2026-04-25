package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditAction represents the type of operation performed.
type AuditAction string

const (
	ActionSet    AuditAction = "set"
	ActionDelete AuditAction = "delete"
	ActionImport AuditAction = "import"
	ActionCopy   AuditAction = "copy"
	ActionRename AuditAction = "rename"
)

// AuditEntry records a single mutation event.
type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Target    string      `json:"target"`
	Action    AuditAction `json:"action"`
	Key       string      `json:"key"`
	Encrypted bool        `json:"encrypted"`
}

// AuditLog manages reading and writing audit entries.
type AuditLog struct {
	path string
}

// NewAuditLog creates an AuditLog backed by a file in the given directory.
func NewAuditLog(dir string) *AuditLog {
	return &AuditLog{path: filepath.Join(dir, "audit.log")}
}

// Append writes a new entry to the audit log.
func (a *AuditLog) Append(entry AuditEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open audit log: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("write audit entry: %w", err)
	}
	return nil
}

// ReadAll returns all entries from the audit log.
func (a *AuditLog) ReadAll() ([]AuditEntry, error) {
	f, err := os.Open(a.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open audit log: %w", err)
	}
	defer f.Close()
	var entries []AuditEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e AuditEntry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode audit entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
