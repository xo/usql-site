package main

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in         string
		cols, rows int
		ok         bool
	}{
		{"", 100, 26, true},
		{"80x24", 80, 24, true},
		{"80", 0, 0, false},
		{"axb", 0, 0, false},
		{"10x24", 0, 0, false},
	}
	for _, test := range tests {
		cols, rows, err := size(test.in)
		if (err == nil) != test.ok || cols != test.cols || rows != test.rows {
			t.Errorf("size(%q) = %d, %d, %v", test.in, cols, rows, err)
		}
	}
}

func TestStripEscapes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"\x1b[32msq:books.db=> \x1b[39m", "sq:books.db=> "},
		{"\x1b[?2004h$ ", "$ "},
		{"\x1b]0;title\x07text", "text"},
		{"\x1b_Gi=31;AAAA\x1b\\done", "done"},
		{"plain", "plain"},
	}
	for _, test := range tests {
		if got := stripEscapes(test.in); got != test.want {
			t.Errorf("stripEscapes(%q) = %q, want %q", test.in, got, test.want)
		}
	}
}

func TestWriteFile(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	vars := map[string]string{"DSN": "pg://localhost/db"}
	if err := writeFile(home, ".config/usql/config.yaml", "connections:\n  pg: ${DSN}\n", vars); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	buf, err := os.ReadFile(filepath.Join(home, ".config", "usql", "config.yaml"))
	if err != nil {
		t.Fatalf("reading the file: %v", err)
	}
	if got, want := string(buf), "connections:\n  pg: pg://localhost/db\n"; got != want {
		t.Errorf("the file holds %q, want %q", got, want)
	}
	if err := writeFile(home, "a.txt", "${OTHER}", vars); err == nil || !strings.Contains(err.Error(), "OTHER") {
		t.Errorf("a variable that is not in needs: got %v", err)
	}
	if err := writeFile(home, "../outside.txt", "x", vars); err == nil {
		t.Error("a file outside the home directory: got no error")
	}
}

func TestCompareVersion(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	usql := filepath.Join(dir, "usql")
	if err := os.WriteFile(usql, []byte("#!/bin/sh\necho 'usql 0.20.0'\n"), 0o755); err != nil {
		t.Fatalf("writing the fake usql: %v", err)
	}
	tests := []struct {
		name, content, want string
	}{
		{"same", "0.20.0\n", ""},
		{"other", "0.19.0\n", "is version 0.20.0, and"},
		{"empty", "\n", "is empty"},
		{"missing", "", "reading the version"},
	}
	for _, test := range tests {
		file := filepath.Join(dir, test.name)
		if test.content != "" {
			if err := os.WriteFile(file, []byte(test.content), 0o644); err != nil {
				t.Fatalf("writing %s: %v", file, err)
			}
		}
		err := compareVersion(t.Context(), usql, file)
		switch {
		case test.want == "" && err != nil:
			t.Errorf("%s: got %v, want no error", test.name, err)
		case test.want != "" && (err == nil || !strings.Contains(err.Error(), test.want)):
			t.Errorf("%s: got %v, want an error with %q", test.name, err, test.want)
		}
	}
}

func TestKeep(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	if err := os.WriteFile(filepath.Join(home, "sales.svg"), []byte("<svg/>"), 0o644); err != nil {
		t.Fatalf("writing the file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, "secret.svg"), []byte("pg://localhost/db"), 0o644); err != nil {
		t.Fatalf("writing the file: %v", err)
	}
	secrets := []string{"pg://localhost/db"}
	kept, err := readKept(home, []string{"sales.svg"}, secrets)
	if err != nil {
		t.Fatalf("readKept: %v", err)
	}
	for _, name := range []string{"missing.svg", "../sales.svg", "secret.svg"} {
		if _, err := readKept(home, []string{name}, secrets); err == nil {
			t.Errorf("readKept(%q): got no error", name)
		}
	}
	dir := filepath.Join(t.TempDir(), "charts", "first")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("making the directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "old.svg"), nil, 0o644); err != nil {
		t.Fatalf("writing the old file: %v", err)
	}
	if err := writeKept(dir, kept); err != nil {
		t.Fatalf("writeKept: %v", err)
	}
	buf, err := os.ReadFile(filepath.Join(dir, "sales.svg"))
	if err != nil || string(buf) != "<svg/>" {
		t.Errorf("the kept file holds %q, %v, want %q", buf, err, "<svg/>")
	}
	if _, err := os.Stat(filepath.Join(dir, "old.svg")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the file of an old recording is still there: %v", err)
	}
}

func TestSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	vars := map[string]string{"PW": "s3cr3t"}
	var buf bytes.Buffer
	if err := do(ctx, &buf, nil, Step{Secret: "PW"}, vars); err != nil {
		t.Fatalf("do: %v", err)
	}
	if got, want := buf.String(), "s3cr3t\r"; got != want {
		t.Errorf("the step typed %q, want %q", got, want)
	}
	buf.Reset()
	if err := do(ctx, &buf, nil, Step{Secret: "OTHER"}, vars); err == nil || !strings.Contains(err.Error(), "OTHER") {
		t.Errorf("a variable that is not in needs: got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("a variable that is not in needs: the step typed %q", buf.String())
	}
}
