package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// apiJSON is a page of the GitHub API, cut down to the fields that the
// program reads and a few that it does not. It holds a draft, a release with
// no name, and line ends of GitHub's editor, and it is not sorted.
const apiJSON = `[
  {
    "tag_name": "v0.2.0",
    "name": "usql v0.2.0",
    "body": "## What's Changed\r\n* A fix\r\n",
    "html_url": "https://github.com/xo/usql/releases/tag/v0.2.0",
    "draft": false,
    "prerelease": false,
    "published_at": "2026-01-08T01:46:41Z"
  },
  {
    "tag_name": "v0.3.0",
    "name": "",
    "body": "",
    "html_url": "https://github.com/xo/usql/releases/tag/v0.3.0",
    "draft": true,
    "published_at": null
  },
  {
    "tag_name": "v0.2.1",
    "name": "  ",
    "body": "  <b>raw</b> ![x](https://example.com/x.png)",
    "html_url": "https://github.com/xo/usql/releases/tag/v0.2.1",
    "draft": false,
    "published_at": "2026-01-08T23:58:02-05:00"
  },
  {
    "tag_name": "v0.1.0",
    "name": "usql v0.1.0",
    "body": null,
    "html_url": "https://github.com/xo/usql/releases/tag/v0.1.0",
    "draft": false,
    "published_at": "2017-03-06T14:47:45Z"
  }
]`

func TestConvertReleases(t *testing.T) {
	t.Parallel()
	var in []apiRelease
	if err := json.Unmarshal([]byte(apiJSON), &in); err != nil {
		t.Fatalf("decoding the API JSON: %v", err)
	}
	want := []release{
		{
			Tag:  "v0.2.1",
			Name: "v0.2.1",
			Date: "2026-01-09",
			Body: "<b>raw</b> ![x](https://example.com/x.png)",
			URL:  "https://github.com/xo/usql/releases/tag/v0.2.1",
		},
		{
			Tag:  "v0.2.0",
			Name: "usql v0.2.0",
			Date: "2026-01-08",
			Body: "## What's Changed\n* A fix",
			URL:  "https://github.com/xo/usql/releases/tag/v0.2.0",
		},
		{
			Tag:  "v0.1.0",
			Name: "usql v0.1.0",
			Date: "2017-03-06",
			Body: "",
			URL:  "https://github.com/xo/usql/releases/tag/v0.1.0",
		},
	}
	if got := convertReleases(in); !reflect.DeepEqual(got, want) {
		t.Errorf("convertReleases =\n%+v\nwant\n%+v", got, want)
	}
	if len(in) != 4 {
		t.Errorf("convertReleases changed its input to %d releases", len(in))
	}
}

func TestEncodeReleasesNone(t *testing.T) {
	t.Parallel()
	if _, err := encodeReleases(nil); err == nil {
		t.Error("encodeReleases(nil) returned no error")
	}
}

func TestNextLink(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{`<https://api.github.com/x?page=2>; rel="next", <https://api.github.com/x?page=4>; rel="last"`, "https://api.github.com/x?page=2"},
		{`<https://api.github.com/x?page=1>; rel="prev", <https://api.github.com/x?page=1>; rel="first"`, ""},
		{`https://api.github.com/x?page=2; rel="next"`, ""},
	}
	for _, test := range tests {
		if got := nextLink(test.in); got != test.want {
			t.Errorf("nextLink(%q) = %q, want %q", test.in, got, test.want)
		}
	}
}

// TestFetchReleases fetches two pages from a server on the loopback address.
// It needs no network.
func TestFetchReleases(t *testing.T) {
	t.Parallel()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Errorf("page %s: Authorization = %q", r.URL.Query().Get("page"), got)
		}
		switch r.URL.Query().Get("page") {
		case "":
			w.Header().Set("Link", fmt.Sprintf(`<%s/releases?page=2>; rel="next"`, srv.URL))
			fmt.Fprint(w, `[{"tag_name": "v2", "published_at": "2026-01-02T00:00:00Z"}]`)
		case "2":
			fmt.Fprint(w, `[{"tag_name": "v1", "published_at": "2026-01-01T00:00:00Z"}]`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	got, err := fetchReleases(context.Background(), srv.Client(), srv.URL+"/releases", "secret")
	if err != nil {
		t.Fatalf("fetchReleases: %v", err)
	}
	var tags []string
	for _, r := range got {
		tags = append(tags, r.TagName)
	}
	if want := []string{"v2", "v1"}; !slices.Equal(tags, want) {
		t.Errorf("fetchReleases returned %v, want %v", tags, want)
	}
}

// TestFetchReleasesOtherHost checks that a next link to another host stops
// the fetch, so that the token goes to no other host.
func TestFetchReleasesOtherHost(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Link", `<http://other.invalid/releases?page=2>; rel="next"`)
		fmt.Fprint(w, `[]`)
	}))
	defer srv.Close()
	if _, err := fetchReleases(context.Background(), srv.Client(), srv.URL+"/releases", "secret"); err == nil {
		t.Error("fetchReleases followed a next link to another host")
	}
}

func TestFetchReleasesStatus(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message": "API rate limit exceeded"}`, http.StatusForbidden)
	}))
	defer srv.Close()
	if _, err := fetchReleases(context.Background(), srv.Client(), srv.URL+"/releases", ""); err == nil {
		t.Error("fetchReleases returned no error for 403 Forbidden")
	}
}

// TestReleasesFile checks that the committed releases.json is sorted newest
// first and that each release has a tag, a date and a URL. It does not fetch
// the releases, so it does not check that the file is current.
func TestReleasesFile(t *testing.T) {
	t.Parallel()
	buf, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(releasesFile)))
	if err != nil {
		t.Fatalf("reading %s: %v. Run go run ./cmd/gen -releases", releasesFile, err)
	}
	var list releaseList
	if err := json.Unmarshal(buf, &list); err != nil {
		t.Fatalf("decoding %s: %v", releasesFile, err)
	}
	if len(list.Releases) == 0 {
		t.Fatalf("%s holds no release", releasesFile)
	}
	for i, r := range list.Releases {
		if r.Tag == "" || r.Date == "" || r.URL == "" {
			t.Errorf("%s: release %d has no tag, date or URL: %+v", releasesFile, i, r)
		}
		if i > 0 && list.Releases[i-1].Date < r.Date {
			t.Errorf("%s: %s comes before %s, which is newer", releasesFile, list.Releases[i-1].Tag, r.Tag)
		}
	}
}
