package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

// releasesFile is the file that the program writes with -releases, from the
// repository root.
const releasesFile = "src/data/generated/releases.json"

// releasesURL is the first page of the releases of usql in the GitHub API.
const releasesURL = "https://api.github.com/repos/xo/usql/releases?per_page=100"

// apiRelease is the part of a release of the GitHub API that the site uses.
type apiRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	Draft   bool   `json:"draft"`
	// PublishedAt is null for a draft.
	PublishedAt *time.Time `json:"published_at"`
}

// releaseList is the content of releases.json.
type releaseList struct {
	// Releases are sorted newest first.
	Releases []release `json:"releases"`
}

// release is one release of usql, as the ::releases directive of the pages
// reads it (src/lib/releases.mjs).
type release struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
	// Date is the day the release was published, in UTC, as 2006-01-02.
	Date string `json:"date"`
	// Body is the Markdown that the maintainers wrote for the release.
	Body string `json:"body"`
	// URL is the page of the release on GitHub.
	URL string `json:"url"`
}

// generateReleases fetches the releases of usql and returns the content of
// releases.json.
func generateReleases(ctx context.Context) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	token := cmp.Or(os.Getenv("GITHUB_TOKEN"), os.Getenv("GH_TOKEN"))
	in, err := fetchReleases(ctx, client, releasesURL, token)
	if err != nil {
		return nil, err
	}
	return encodeReleases(convertReleases(in))
}

// fetchReleases returns the releases of every page of the GitHub API that
// starts at first. It follows the next link of each page. It sends token,
// when it is not empty, only to the host of first.
func fetchReleases(ctx context.Context, client *http.Client, first, token string) ([]apiRelease, error) {
	start, err := url.Parse(first)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", first, err)
	}
	var all []apiRelease
	for next := first; next != ""; {
		u, err := url.Parse(next)
		if err != nil {
			return nil, fmt.Errorf("parsing the next page %s: %w", next, err)
		}
		if u.Scheme != start.Scheme || u.Host != start.Host {
			return nil, fmt.Errorf("the next page %s is on another host than %s", next, first)
		}
		page, link, err := fetchPage(ctx, client, next, token)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		next = nextLink(link)
	}
	return all, nil
}

// fetchPage returns the releases of one page of the GitHub API, and its Link
// header.
func fetchPage(ctx context.Context, client *http.Client, pageURL, token string) ([]apiRelease, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("making the request for %s: %w", pageURL, err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "usql-site-gen")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("getting %s: %w", pageURL, err)
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading %s: %w", pageURL, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("getting %s: %s: %s", pageURL, res.Status, bytes.TrimSpace(buf))
	}
	var page []apiRelease
	if err := json.Unmarshal(buf, &page); err != nil {
		return nil, "", fmt.Errorf("decoding %s: %w", pageURL, err)
	}
	return page, res.Header.Get("Link"), nil
}

// nextLink returns the URL with rel="next" in the Link header link, or an
// empty string when there is none.
func nextLink(link string) string {
	for part := range strings.SplitSeq(link, ",") {
		target, params, ok := strings.Cut(part, ";")
		if !ok {
			continue
		}
		target = strings.TrimSpace(target)
		if !strings.HasPrefix(target, "<") || !strings.HasSuffix(target, ">") {
			continue
		}
		for param := range strings.SplitSeq(params, ";") {
			if strings.TrimSpace(param) == `rel="next"` {
				return target[1 : len(target)-1]
			}
		}
	}
	return ""
}

// convertReleases returns the published releases of in, newest first. It
// leaves out the drafts. A body has line feeds for line ends, and no space at
// its start or end. A release with no name takes its tag as its name.
func convertReleases(in []apiRelease) []release {
	published := slices.DeleteFunc(slices.Clone(in), func(r apiRelease) bool {
		return r.Draft || r.PublishedAt == nil
	})
	slices.SortStableFunc(published, func(a, b apiRelease) int {
		return b.PublishedAt.Compare(*a.PublishedAt)
	})
	out := make([]release, 0, len(published))
	for _, r := range published {
		body := strings.ReplaceAll(r.Body, "\r\n", "\n")
		out = append(out, release{
			Tag:  r.TagName,
			Name: cmp.Or(strings.TrimSpace(r.Name), r.TagName),
			Date: r.PublishedAt.UTC().Format(time.DateOnly),
			Body: strings.TrimSpace(body),
			URL:  r.HTMLURL,
		})
	}
	return out
}

// encodeReleases returns the content of releases.json.
func encodeReleases(releases []release) ([]byte, error) {
	if len(releases) == 0 {
		return nil, errors.New("the GitHub API returned no published release")
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	if err := enc.Encode(releaseList{Releases: releases}); err != nil {
		return nil, fmt.Errorf("encoding the releases: %w", err)
	}
	return buf.Bytes(), nil
}
