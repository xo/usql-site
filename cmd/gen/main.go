// Command gen writes the data files of the site that come from a Go module
// that go.mod pins. It writes src/data/generated/schemes.json, the URL
// schemes of dburl's registry. The ::schemes and ::databases directives of
// the pages read that file (src/lib/schemes.mjs). Run it in the root of the
// repository after go.mod moves to another release of dburl:
//
//	go run ./cmd/gen
//
// A test fails when the file differs from what the program makes. See D7.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/xo/dburl"
)

// schemesFile is the file that the program writes, from the repository root.
const schemesFile = "src/data/generated/schemes.json"

func main() {
	out := flag.String("out", schemesFile, "the file to write")
	flag.Parse()
	if err := run(*out); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

// run writes the schemes to out.
func run(out string) error {
	if _, err := os.Stat("astro.config.mjs"); err != nil {
		return errors.New("run it in the root of the repository")
	}
	buf, err := generate()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("making the directory of %s: %w", out, err)
	}
	if err := os.WriteFile(out, buf, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", out, err)
	}
	return nil
}

// registry is the content of schemes.json.
type registry struct {
	// Dburl is the release of dburl that go.mod pins.
	Dburl   string   `json:"dburl"`
	Schemes []scheme `json:"schemes"`
}

// scheme is one scheme of dburl's registry. The fields have the meaning of
// the fields of dburl.Scheme.
type scheme struct {
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Dialect string `json:"dialect"`
	// Aliases are the names that open the scheme, other than Name. They hold
	// the short alias that dburl adds, in dburl's order: the shortest first.
	Aliases []string `json:"aliases"`
	// Opaque is true when the rest of the URL is a path, as in sq:books.db.
	Opaque    bool     `json:"opaque"`
	Transport []string `json:"transport"`
	// FileType is true when usql opens a file of this database from its path
	// alone, and with the scheme file:.
	FileType   bool     `json:"file_type"`
	Driver     string   `json:"driver"`
	DriverURL  string   `json:"driver_url"`
	Home       string   `json:"home"`
	CGO        bool     `json:"cgo"`
	Deployment []string `json:"deployment"`
}

// generate returns the content of schemes.json.
func generate() ([]byte, error) {
	version, err := dburlVersion()
	if err != nil {
		return nil, err
	}
	fileTypes := dburl.FileTypes()
	var schemes []scheme
	for _, s := range dburl.BaseSchemes() {
		_, aliases := dburl.SchemeNameAndAliases(s.Name)
		schemes = append(schemes, scheme{
			Name:       s.Name,
			Desc:       s.Desc,
			Dialect:    s.Dialect,
			Aliases:    nonNil(aliases),
			Opaque:     s.Opaque,
			Transport:  transports(s.Transport),
			FileType:   slices.Contains(fileTypes, s.Name),
			Driver:     s.GoPackage,
			DriverURL:  s.DriverURL,
			Home:       s.Home,
			CGO:        s.RequiresCGO,
			Deployment: deployments(s.Deployment),
		})
	}
	slices.SortFunc(schemes, func(a, b scheme) int { return strings.Compare(a.Name, b.Name) })
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	if err := enc.Encode(registry{Dburl: version, Schemes: schemes}); err != nil {
		return nil, fmt.Errorf("encoding the schemes: %w", err)
	}
	return buf.Bytes(), nil
}

// dburlVersion returns the release of dburl that the program is built with,
// which is the one that go.mod pins.
func dburlVersion() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("the program has no build information")
	}
	for _, dep := range info.Deps {
		if dep.Path == "github.com/xo/dburl" {
			if dep.Replace != nil {
				return "", errors.New("go.mod replaces github.com/xo/dburl. Remove the replace directive")
			}
			return dep.Version, nil
		}
	}
	return "", errors.New("the build information does not name github.com/xo/dburl")
}

// transports returns the names of the transports of t, as a URL writes
// them after the +.
func transports(t dburl.Transport) []string {
	names := []string{}
	for _, v := range []struct {
		t    dburl.Transport
		name string
	}{
		{dburl.TransportTCP, "tcp"},
		{dburl.TransportUDP, "udp"},
		{dburl.TransportUnix, "unix"},
		{dburl.TransportAny, "any"},
	} {
		if t&v.t != 0 {
			names = append(names, v.name)
		}
	}
	return names
}

// deployments returns the names of the deployment kinds of d.
func deployments(d dburl.Deployment) []string {
	names := []string{}
	for _, v := range []struct {
		d    dburl.Deployment
		name string
	}{
		{dburl.DeploymentEmbedded, "embedded"},
		{dburl.DeploymentServer, "server"},
		{dburl.DeploymentHosted, "hosted"},
	} {
		if d&v.d != 0 {
			names = append(names, v.name)
		}
	}
	return names
}

// nonNil returns s, or an empty slice when s is nil, so that the JSON holds
// [] and not null.
func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
