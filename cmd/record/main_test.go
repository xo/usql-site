package main

import (
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
