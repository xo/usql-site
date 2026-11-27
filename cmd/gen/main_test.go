package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xo/dburl"
)

// TestSchemesAreCurrent checks that the committed schemes.json is what the
// program makes from the release of dburl that go.mod pins.
func TestSchemesAreCurrent(t *testing.T) {
	t.Parallel()
	want, err := generate()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	got, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(schemesFile)))
	if err != nil {
		t.Fatalf("reading %s: %v. Run go run ./cmd/gen", schemesFile, err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s differs from what cmd/gen makes from dburl. Run go run ./cmd/gen", schemesFile)
	}
}

func TestTransports(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   dburl.Transport
		want string
	}{
		{dburl.TransportNone, ""},
		{dburl.TransportTCP, "tcp"},
		{dburl.TransportTCP | dburl.TransportUDP | dburl.TransportUnix, "tcp udp unix"},
		{dburl.TransportAny, "any"},
	}
	for _, test := range tests {
		if got := strings.Join(transports(test.in), " "); got != test.want {
			t.Errorf("transports(%d) = %q, want %q", test.in, got, test.want)
		}
	}
}
