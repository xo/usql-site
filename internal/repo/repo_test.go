// Package repo holds the tests of the repository itself: the files that D1
// requires, the decision index, and the recordings that the pages name.
package repo

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// root is the repository root, from the directory of this package.
var root = filepath.Join("..", "..")

// TestClaudeImportsAgents checks that CLAUDE.md holds only the import of
// AGENTS.md, so that every agent reads the same rules. See D1.
func TestClaudeImportsAgents(t *testing.T) {
	t.Parallel()
	if got := strings.TrimSpace(readFile(t, "CLAUDE.md")); got != "@AGENTS.md" {
		t.Errorf("CLAUDE.md holds %q, and must hold only @AGENTS.md. Put the rules in AGENTS.md", got)
	}
}

// TestTheRootHoldsFourDocuments checks that the only Markdown files in the
// root are the four that D1 names.
func TestTheRootHoldsFourDocuments(t *testing.T) {
	t.Parallel()
	want := []string{"AGENTS.md", "CLAUDE.md", "CONTRIBUTING.md", "README.md"}
	matches, err := filepath.Glob(filepath.Join(root, "*.md"))
	if err != nil {
		t.Fatalf("listing the root: %v", err)
	}
	got := make([]string, 0, len(matches))
	for _, m := range matches {
		got = append(got, filepath.Base(m))
	}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Errorf("the root holds %v, and must hold %v. Move any other document to docs/", got, want)
	}
}

// TestEveryDocumentIsInTheTable checks that AGENTS.md names every document
// of docs/, and that every document it names exists.
func TestEveryDocumentIsInTheTable(t *testing.T) {
	t.Parallel()
	agents := readFile(t, "AGENTS.md")
	matches, err := filepath.Glob(filepath.Join(root, "docs", "*.md"))
	if err != nil {
		t.Fatalf("listing docs: %v", err)
	}
	for _, m := range matches {
		name := "docs/" + filepath.Base(m)
		if !strings.Contains(agents, "("+name+")") {
			t.Errorf("%s is not in the table of AGENTS.md. Add a row for it, or delete it", name)
		}
	}
	for _, m := range regexp.MustCompile(`\]\((docs/[^)#]+)\)`).FindAllStringSubmatch(agents, -1) {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(m[1]))); err != nil {
			t.Errorf("AGENTS.md names %s, which does not exist", m[1])
		}
	}
}

// TestSkillsAreCopies checks that each skill of skills-lock.json is an
// ordinary folder in .agents/skills and in .claude/skills, and that the two
// copies are the same. See D1.
func TestSkillsAreCopies(t *testing.T) {
	t.Parallel()
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal([]byte(readFile(t, "skills-lock.json")), &lock); err != nil {
		t.Fatalf("reading skills-lock.json: %v", err)
	}
	if len(lock.Skills) == 0 {
		t.Fatal("skills-lock.json names no skill")
	}
	dirs := []string{filepath.Join(".agents", "skills"), filepath.Join(".claude", "skills")}
	for _, dir := range dirs {
		entries, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			t.Errorf("%s: %v. Install the skills with the command in CONTRIBUTING.md", dir, err)
			continue
		}
		for _, e := range entries {
			if _, ok := lock.Skills[e.Name()]; !ok {
				t.Errorf("%s is not in skills-lock.json. Add it with the command in CONTRIBUTING.md",
					filepath.Join(dir, e.Name()))
			}
		}
	}
	for _, name := range slices.Sorted(maps.Keys(lock.Skills)) {
		agents := skillFiles(t, filepath.Join(dirs[0], name))
		claude := skillFiles(t, filepath.Join(dirs[1], name))
		if agents == nil || claude == nil {
			continue
		}
		if !maps.Equal(agents, claude) {
			t.Errorf("%s: the copies in %s and %s differ. Install the skill again with --copy",
				name, dirs[0], dirs[1])
		}
	}
}

// skillFiles returns the content of every file of one copy of a skill, by
// its path in the copy. It returns nil after it reports a copy that is
// missing or that is a symbolic link.
func skillFiles(t *testing.T, dir string) map[string]string {
	t.Helper()
	path := filepath.Join(root, dir)
	fi, err := os.Lstat(path)
	switch {
	case err != nil:
		t.Errorf("%s is missing. Install the skill with the command in CONTRIBUTING.md", dir)
		return nil
	case !fi.IsDir():
		t.Errorf("%s is not a folder. Install the skill with --copy. See D1", dir)
		return nil
	}
	files := make(map[string]string)
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		buf, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = string(buf)
		return nil
	})
	if err != nil {
		t.Errorf("reading %s: %v", dir, err)
		return nil
	}
	return files
}

// decision is one file of docs/decisions.
type decision struct {
	num, file, title, status string
}

// decisions reads every decision of docs/decisions.
func decisions(t *testing.T) []decision {
	t.Helper()
	dir := filepath.Join(root, "docs", "decisions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading docs/decisions: %v", err)
	}
	nameRE := regexp.MustCompile(`^D(\d{3})-[a-z0-9-]+\.md$`)
	headRE := regexp.MustCompile(`\A# D(\d+)\. (.+)\n\nStatus: (.+?)\.?\n`)
	var ds []decision
	for _, e := range entries {
		if e.Name() == "README.md" {
			continue
		}
		m := nameRE.FindStringSubmatch(e.Name())
		if m == nil {
			t.Errorf("docs/decisions/%s: the name must be D<nnn>-<title>.md. See D1", e.Name())
			continue
		}
		head := headRE.FindStringSubmatch(readFile(t, filepath.Join("docs", "decisions", e.Name())))
		if head == nil {
			t.Errorf("docs/decisions/%s must open with \"# Dn. Title\", a blank line, and \"Status: ...\"", e.Name())
			continue
		}
		num := strings.TrimLeft(m[1], "0")
		if head[1] != num {
			t.Errorf("docs/decisions/%s opens with D%s", e.Name(), head[1])
		}
		ds = append(ds, decision{num: num, file: e.Name(), title: head[2], status: head[3]})
	}
	return ds
}

// TestTheDecisionIndexIsComplete checks that every decision has a row in the
// index, with its file, its title and its status, and that every row has a
// decision.
func TestTheDecisionIndexIsComplete(t *testing.T) {
	t.Parallel()
	index := readFile(t, filepath.Join("docs", "decisions", "README.md"))
	// a row is compared by its cells, so that the padding of a formatter does
	// not matter
	rows := make(map[string]string)
	for _, m := range regexp.MustCompile(`(?m)^\|\s*\[D(\d+)\]\(.*$`).FindAllStringSubmatch(index, -1) {
		rows[m[1]] = row(m[0])
	}
	written := make(map[string]bool)
	for _, d := range decisions(t) {
		written[d.num] = true
		want := fmt.Sprintf("| [D%s](%s) | %s | %s |", d.num, d.file, d.title, d.status)
		switch got, ok := rows[d.num]; {
		case !ok:
			t.Errorf("D%s has no row in docs/decisions/README.md. Add:\n%s", d.num, want)
		case got != want:
			t.Errorf("D%s: the row in docs/decisions/README.md is\n%s\nand the file says\n%s", d.num, got, want)
		}
	}
	for _, num := range slices.Sorted(maps.Keys(rows)) {
		if !written[num] {
			t.Errorf("docs/decisions/README.md has a row for D%s, and no file holds it", num)
		}
	}
}

// row returns a row of a Markdown table with one space around each cell.
func row(s string) string {
	cells := strings.Split(strings.Trim(strings.TrimSpace(s), "|"), "|")
	for i, c := range cells {
		cells[i] = strings.TrimSpace(c)
	}
	return "| " + strings.Join(cells, " | ") + " |"
}

// TestEveryRecordingHasAScript checks that every recording of public/casts
// has its script in recordings, and every script its recording, so that
// every recording can be made again. See D5.
func TestEveryRecordingHasAScript(t *testing.T) {
	t.Parallel()
	casts := names(t, filepath.Join("public", "casts"), ".cast")
	scripts := names(t, "recordings", ".yml")
	for _, name := range casts {
		if !slices.Contains(scripts, name) {
			t.Errorf("public/casts/%s.cast has no script recordings/%s.yml. Write one, and record it with cmd/record", name, name)
		}
	}
	for _, name := range scripts {
		if !slices.Contains(casts, name) {
			t.Errorf("recordings/%s.yml has no recording. Run go run ./cmd/record %s", name, name)
		}
	}
}

// TestEveryPageIsMarkdown checks that no page is MDX, which the editor cannot
// open. See D4.
func TestEveryPageIsMarkdown(t *testing.T) {
	t.Parallel()
	err := filepath.WalkDir(filepath.Join(root, "src", "content", "docs"), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Ext(p) != ".md" {
			rel, _ := filepath.Rel(root, p)
			t.Errorf("%s is not a .md file. A page is Markdown, so that the editor can open it. See D4", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("reading the pages: %v", err)
	}
}

// names returns the names of the files of dir with the extension ext,
// without it.
func names(t *testing.T, dir, ext string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(root, dir, "*"+ext))
	if err != nil {
		t.Fatalf("listing %s: %v", dir, err)
	}
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, strings.TrimSuffix(filepath.Base(m), ext))
	}
	return out
}

// readFile returns the content of the file name of the repository.
func readFile(t *testing.T, name string) string {
	t.Helper()
	buf, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return string(buf)
}
