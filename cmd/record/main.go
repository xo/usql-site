// Command record makes the terminal recordings of the site. It reads a script
// from recordings/<name>.yml, runs a shell in a pseudo-terminal, types each
// step of the script into it, and writes what the terminal showed to
// public/casts/<name>.cast in the asciicast v2 format, which the player of the
// site reads.
//
//	go run ./cmd/record tour
//	go run ./cmd/record -all
//
// A script runs in a new temporary directory, which is also the home
// directory of the shell, so that no configuration of the person who records
// it reaches the recording. See docs/RECORDING.md.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/goccy/go-yaml"
)

func main() {
	all := flag.Bool("all", false, "record every script in the scripts directory")
	dir := flag.String("dir", "recordings", "the directory of the scripts")
	out := flag.String("out", filepath.Join("public", "casts"), "the directory of the recordings")
	flag.Parse()
	names := flag.Args()
	if *all {
		var err error
		if names, err = scripts(*dir); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "usage: record [-all] [name...]")
		os.Exit(2)
	}
	ctx := context.Background()
	failed := false
	for _, name := range names {
		if err := record(ctx, filepath.Join(*dir, name+".yml"), filepath.Join(*out, name+".cast")); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", name, err)
			failed = true
			continue
		}
		fmt.Println("recorded", name)
	}
	if failed {
		os.Exit(1)
	}
}

// scripts returns the names of the scripts in dir.
func scripts(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, fmt.Errorf("listing %s: %w", dir, err)
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(file), ".yml"))
	}
	slices.Sort(names)
	return names, nil
}

// Script is one recording.
type Script struct {
	// Title is the title of the recording.
	Title string `yaml:"title"`
	// Size is the size of the terminal, as COLSxROWS.
	Size string `yaml:"size"`
	// Needs names the variables of the environment that the script needs,
	// such as the address of a database. A script whose variable is not set
	// is skipped with an error.
	Needs []string `yaml:"needs"`
	// Files are written to the home directory before the recording, with
	// ${NAME} replaced by the variable NAME of Needs.
	Files map[string]string `yaml:"files"`
	// Setup runs in the home directory before the recording, one command at
	// a time, with sh -c. It is not recorded.
	Setup []string `yaml:"setup"`
	// Steps are typed into the terminal, in order.
	Steps []Step `yaml:"steps"`
}

// Step is one step of a script. Exactly one of its fields is set.
type Step struct {
	// Type types the text and then Enter.
	Type string `yaml:"type"`
	// Keys types the text with no Enter. Escape sequences such as \t and \x1b
	// are read as in a Go string.
	Keys string `yaml:"keys"`
	// Wait waits until the terminal shows the text.
	Wait string `yaml:"wait"`
	// Sleep pauses, as a Go duration such as 1.5s.
	Sleep string `yaml:"sleep"`
}

// The timing of the typing and of the waits.
const (
	keyDelay    = 45 * time.Millisecond
	enterDelay  = 300 * time.Millisecond
	waitTimeout = 20 * time.Second
	exitTimeout = 10 * time.Second
	// idleLimit is the longest pause that the player shows.
	idleLimit = 2.0
)

// record runs the script in file and writes the recording to out.
func record(ctx context.Context, file, out string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("reading script: %w", err)
	}
	var s Script
	if err := yaml.UnmarshalWithOptions(buf, &s, yaml.Strict()); err != nil {
		return fmt.Errorf("parsing %s: %w", file, err)
	}
	cols, rows, err := size(s.Size)
	if err != nil {
		return err
	}
	vars := make(map[string]string, len(s.Needs))
	for _, name := range s.Needs {
		v, ok := os.LookupEnv(name)
		if !ok || v == "" {
			return fmt.Errorf("the script needs the variable %s", name)
		}
		vars[name] = v
	}
	home, err := os.MkdirTemp("", "usql-record-")
	if err != nil {
		return fmt.Errorf("making the home directory: %w", err)
	}
	defer os.RemoveAll(home)
	env := environ(home)
	for name, content := range s.Files {
		if err := writeFile(home, name, content, vars); err != nil {
			return err
		}
	}
	for _, command := range s.Setup {
		cmd := exec.CommandContext(ctx, "sh", "-c", command)
		cmd.Dir, cmd.Env = home, env
		if buf, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("running setup %q: %w: %s", command, err, buf)
		}
	}
	events, err := run(ctx, home, env, cols, rows, s.Steps, slices.Collect(maps.Values(vars)))
	if err != nil {
		return err
	}
	return writeCast(out, s.Title, cols, rows, events)
}

// size parses COLSxROWS.
func size(s string) (int, int, error) {
	if s == "" {
		return 100, 26, nil
	}
	c, r, ok := strings.Cut(s, "x")
	cols, err1 := strconv.Atoi(c)
	rows, err2 := strconv.Atoi(r)
	if !ok || err1 != nil || err2 != nil || cols < 20 || rows < 5 {
		return 0, 0, fmt.Errorf("invalid size %q, want COLSxROWS", s)
	}
	return cols, rows, nil
}

// environ returns the environment of the shell. Only PATH comes from the
// environment of the recorder, so that the recording is the same on every
// machine.
func environ(home string) []string {
	return []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + home,
		"XDG_CONFIG_HOME=" + filepath.Join(home, ".config"),
		"XDG_STATE_HOME=" + filepath.Join(home, ".local", "state"),
		"XDG_CACHE_HOME=" + filepath.Join(home, ".cache"),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"LANG=C.UTF-8",
		"PS1=$ ",
		// the player draws no images, so usql asks the terminal nothing
		"TERM_GRAPHICS=none",
	}
}

// writeFile writes content to name under home, with ${NAME} replaced by the
// variable NAME.
func writeFile(home, name, content string, vars map[string]string) error {
	path := filepath.Join(home, filepath.FromSlash(name))
	if !strings.HasPrefix(path, home+string(filepath.Separator)) {
		return fmt.Errorf("the file %s is outside the home directory", name)
	}
	var missing []string
	content = os.Expand(content, func(key string) string {
		v, ok := vars[key]
		if !ok {
			missing = append(missing, key)
		}
		return v
	})
	if len(missing) != 0 {
		return fmt.Errorf("the file %s names %s, which is not in needs", name, strings.Join(missing, ", "))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("making the directory of %s: %w", name, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", name, err)
	}
	return nil
}

// event is one event of a recording: the time in seconds and the output.
type event struct {
	t    float64
	data string
}

// terminal is the output of the shell, as it arrives.
type terminal struct {
	start time.Time

	mu     sync.Mutex
	events []event
	// seen is the output since the last wait, with no escape sequences.
	seen strings.Builder
	// changed is closed and replaced when output arrives.
	changed chan struct{}
}

// read copies the output of r into t until r ends.
func (t *terminal) read(r io.Reader) {
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			t.mu.Lock()
			t.events = append(t.events, event{time.Since(t.start).Seconds(), string(buf[:n])})
			t.seen.WriteString(stripEscapes(string(buf[:n])))
			close(t.changed)
			t.changed = make(chan struct{})
			t.mu.Unlock()
		}
		if err != nil {
			return
		}
	}
}

// wait waits until the output since the last wait holds s.
func (t *terminal) wait(ctx context.Context, s string) error {
	ctx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()
	for {
		t.mu.Lock()
		seen := t.seen.String()
		changed := t.changed
		if i := strings.Index(seen, s); i != -1 {
			rest := seen[i+len(s):]
			t.seen.Reset()
			t.seen.WriteString(rest)
			t.mu.Unlock()
			return nil
		}
		t.mu.Unlock()
		select {
		case <-changed:
		case <-ctx.Done():
			return fmt.Errorf("waiting for %q: the terminal showed %q", s, tail(seen, 200))
		}
	}
}

// run runs the steps in a shell, and returns the events of the recording.
// It fails when the recording shows one of secrets.
func run(ctx context.Context, home string, env []string, cols, rows int, steps []Step, secrets []string) ([]event, error) {
	cmd := exec.CommandContext(ctx, "bash", "--noprofile", "--norc")
	cmd.Dir, cmd.Env = home, env
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
	if err != nil {
		return nil, fmt.Errorf("starting the shell: %w", err)
	}
	defer f.Close()
	t := &terminal{start: time.Now(), changed: make(chan struct{})}
	done := make(chan struct{})
	go func() {
		t.read(f)
		close(done)
	}()
	if err := t.wait(ctx, "$ "); err != nil {
		return nil, err
	}
	for i, step := range steps {
		if err := do(ctx, f, t, step); err != nil {
			return nil, fmt.Errorf("step %d: %w", i+1, err)
		}
	}
	// the recording ends before the exit that closes the shell
	t.mu.Lock()
	end := len(t.events)
	t.mu.Unlock()
	if _, err := io.WriteString(f, "exit\r"); err != nil {
		return nil, fmt.Errorf("ending the shell: %w", err)
	}
	select {
	case <-done:
	case <-time.After(exitTimeout):
		return nil, errors.New("the shell did not exit")
	}
	_ = cmd.Wait()
	t.mu.Lock()
	defer t.mu.Unlock()
	events := t.events[:end]
	for _, e := range events {
		for _, secret := range secrets {
			if strings.Contains(e.data, secret) {
				return nil, errors.New("the recording shows the value of a variable of needs")
			}
		}
	}
	return events, nil
}

// do runs one step.
func do(ctx context.Context, w io.Writer, t *terminal, step Step) error {
	switch {
	case step.Type != "":
		if err := typeText(ctx, w, step.Type); err != nil {
			return err
		}
		time.Sleep(enterDelay)
		_, err := io.WriteString(w, "\r")
		return err
	case step.Keys != "":
		s, err := strconv.Unquote(`"` + strings.ReplaceAll(step.Keys, `"`, `\"`) + `"`)
		if err != nil {
			return fmt.Errorf("reading keys %q: %w", step.Keys, err)
		}
		return typeText(ctx, w, s)
	case step.Wait != "":
		return t.wait(ctx, step.Wait)
	case step.Sleep != "":
		d, err := time.ParseDuration(step.Sleep)
		if err != nil {
			return fmt.Errorf("reading sleep %q: %w", step.Sleep, err)
		}
		time.Sleep(d)
		return nil
	}
	return errors.New("the step is empty")
}

// typeText types s one character at a time.
func typeText(ctx context.Context, w io.Writer, s string) error {
	for _, r := range s {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, err := io.WriteString(w, string(r)); err != nil {
			return fmt.Errorf("typing: %w", err)
		}
		time.Sleep(keyDelay)
	}
	return nil
}

// writeCast writes the events to out in the asciicast v2 format.
func writeCast(out, title string, cols, rows int, events []event) error {
	var buf bytes.Buffer
	header := map[string]any{
		"version":         2,
		"width":           cols,
		"height":          rows,
		"idle_time_limit": idleLimit,
		"env":             map[string]string{"TERM": "xterm-256color", "SHELL": "/bin/bash"},
	}
	if title != "" {
		header["title"] = title
	}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(header); err != nil {
		return fmt.Errorf("encoding the header: %w", err)
	}
	for _, e := range events {
		t := strconv.FormatFloat(e.t, 'f', 3, 64)
		data, err := json.Marshal(e.data)
		if err != nil {
			return fmt.Errorf("encoding an event: %w", err)
		}
		fmt.Fprintf(&buf, "[%s, \"o\", %s]\n", t, data)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("making the directory of %s: %w", out, err)
	}
	if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", out, err)
	}
	return nil
}

// escapeRE matches the escape sequences of a terminal: CSI, OSC and the
// two-byte forms.
var escapeRE = regexp.MustCompile(`\x1b(\[[0-?]*[ -/]*[@-~]|\][^\x07\x1b]*(\x07|\x1b\\)|_[^\x1b]*\x1b\\|[@-Z\\-_])`)

// stripEscapes removes the escape sequences from s.
func stripEscapes(s string) string {
	return escapeRE.ReplaceAllString(s, "")
}

// tail returns the last n bytes of s.
func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
