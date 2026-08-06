// Command download-fixtures regenerates the test fixtures from the live API.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// listFixturePageSize is how many results each list fixture page keeps. The
// list tests assert on every record, so a live page is far too large.
const listFixturePageSize = 2

// maxFetchAttempts bounds the retries on 429. A full run is dozens of requests
// and the sustained window resets daily, so a wait past this is not worth it.
const maxFetchAttempts = 4

// window is the last observed state of one rate-limit window, from Metron's
// X-RateLimit-* headers.
type window struct {
	remaining int
	reset     time.Time
	seen      bool
}

// wait returns how long until w allows another request.
func (w window) wait() time.Duration {
	if !w.seen || w.remaining > 0 {
		return 0
	}

	return max(0, time.Until(w.reset))
}

// observe records one window's headers, leaving prior state if neither is set.
func (w *window) observe(h http.Header, remainingHeader, resetHeader string) {
	remaining, errRemaining := strconv.Atoi(h.Get(remainingHeader))
	reset, errReset := strconv.ParseInt(h.Get(resetHeader), 10, 64)

	if errRemaining != nil || errReset != nil {
		return
	}

	w.remaining = remaining
	w.reset = time.Unix(reset, 0)
	w.seen = true
}

// limiter paces the run off Metron's rate-limit headers, mirroring the reactive
// approach in the library itself: the sustained limit varies per account tier,
// so an exhausted window is only ever known from a response.
type limiter struct {
	burst     window
	sustained window
}

func (l *limiter) throttle() {
	if d := max(l.burst.wait(), l.sustained.wait()); d > 0 {
		log.Printf("rate limit exhausted, sleeping %s", d.Truncate(time.Millisecond))
		time.Sleep(d)
	}
}

func (l *limiter) observe(h http.Header) {
	l.burst.observe(h, "X-RateLimit-Burst-Remaining", "X-RateLimit-Burst-Reset")
	l.sustained.observe(h, "X-RateLimit-Sustained-Remaining", "X-RateLimit-Sustained-Reset")
}

// listPage is the pagination envelope. Results stay raw so trimming never
// reserializes a record.
type listPage struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []json.RawMessage `json:"results"`
}

func main() {
	apiToken := flag.String("api-token", "", "api-token")

	flag.Parse()
	c := &http.Client{Timeout: time.Second * 10}

	var l limiter

	for _, v := range []struct {
		URL      string
		FileName string
	}{
		{"https://metron.cloud/api/arc/659/", "arc_659.json"},
		{"https://metron.cloud/api/arc/1493/", "arc_1493.json"},
		{"https://metron.cloud/api/character/83/", "character_83.json"},
		{"https://metron.cloud/api/character/26153/", "character_26153.json"},
		{"https://metron.cloud/api/creator/5958/", "creator_5958.json"},
		{"https://metron.cloud/api/creator/11237/", "creator_11237.json"},
		{"https://metron.cloud/api/imprint/1/", "imprint_1.json"},
		{"https://metron.cloud/api/issue/2558/", "issue_2558.json"},
		{"https://metron.cloud/api/issue/112901/", "issue_112901.json"},
		{"https://metron.cloud/api/publisher/1/", "publisher_1.json"},
		{"https://metron.cloud/api/publisher/29/", "publisher_29.json"},
		{"https://metron.cloud/api/series/793/", "series_793.json"},
		{"https://metron.cloud/api/series/3371/", "series_3371.json"},
		{"https://metron.cloud/api/team/180/", "team_180.json"},
		{"https://metron.cloud/api/team/930/", "team_930.json"},
		{"https://metron.cloud/api/universe/24/", "universe_24.json"},
	} {
		if err := makeFixture(c, *apiToken, v.URL, v.FileName, &l); err != nil {
			log.Println(err)
		}
		time.Sleep(2 * time.Second)
	}

	for _, v := range []struct {
		URL    string
		Prefix string
	}{
		{"https://metron.cloud/api/arc/", "arc_list"},
		{"https://metron.cloud/api/arc/659/issue_list/", "arc_659_issue_list"},
		{"https://metron.cloud/api/character/", "character_list"},
		{"https://metron.cloud/api/character/83/issue_list/", "character_83_issue_list"},
		{"https://metron.cloud/api/creator/", "creator_list"},
		{"https://metron.cloud/api/imprint/", "imprint_list"},
		{"https://metron.cloud/api/issue/", "issue_list"},
		{"https://metron.cloud/api/publisher/", "publisher_list"},
		{"https://metron.cloud/api/publisher/1/series_list/", "publisher_1_series_list"},
		{"https://metron.cloud/api/role/", "role_list"},
		{"https://metron.cloud/api/series/", "series_list"},
		{"https://metron.cloud/api/series/793/issue_list/", "series_793_issue_list"},
		{"https://metron.cloud/api/series_type/", "series_type_list"},
		{"https://metron.cloud/api/team/", "team_list"},
		{"https://metron.cloud/api/team/180/issue_list/", "team_180_issue_list"},
		{"https://metron.cloud/api/universe/", "universe_list"},
	} {
		if err := makeListFixture(c, *apiToken, v.URL, v.Prefix, &l); err != nil {
			log.Println(err)
		}
		time.Sleep(2 * time.Second)
	}
}

// makeListFixture writes the two-page fixture pair the list tests iterate.
// Both pages are cut from one live page: role and series_type hold too few
// records to have a second, and the tests only exercise the pagination
// mechanics. The envelope is rewritten to match, so page two ends iteration.
func makeListFixture(c *http.Client, apiToken, url, prefix string, l *limiter) error {
	b, err := fetch(c, apiToken, url+"?page=1", l)
	if err != nil {
		return err
	}

	var live listPage

	if err = json.Unmarshal(b, &live); err != nil {
		return err
	}

	if len(live.Results) < 2*listFixturePageSize {
		return fmt.Errorf("%s: got %d results, want %d", url, len(live.Results), 2*listFixturePageSize)
	}

	pages := [2]listPage{
		{Next: new(url + "?page=2"), Results: live.Results[:listFixturePageSize]},
		{Previous: new(url + "?page=1"), Results: live.Results[listFixturePageSize : 2*listFixturePageSize]},
	}

	for i, page := range pages {
		page.Count = 2 * listFixturePageSize

		if b, err = marshalFixture(page); err != nil {
			return err
		}

		name := filepath.Join("fixtures/", fmt.Sprintf("%s_%d.json", prefix, i+1))

		if err = os.WriteFile(name, b, 0o600); err != nil {
			return err
		}
	}

	return nil
}

// marshalFixture matches the raw bytes the detail fixtures keep: HTML escaping
// off, and no final newline (see .editorconfig).
func marshalFixture(v any) ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func makeFixture(c *http.Client, apiToken, url, fileName string, l *limiter) error {
	b, err := fetch(c, apiToken, url, l)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer

	if err = json.Indent(&prettyJSON, b, "", "  "); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("fixtures/", fileName), prettyJSON.Bytes(), 0o600)
}

// fetch performs a GET, pacing on the rate-limit headers and retrying a 429 for
// as long as its Retry-After asks.
func fetch(c *http.Client, apiToken, url string, l *limiter) ([]byte, error) {
	for attempt := range maxFetchAttempts {
		l.throttle()

		b, wait, err := tryFetch(c, apiToken, url, l)
		if err != nil {
			return nil, err
		}

		if wait == 0 {
			return b, nil
		}

		if attempt == maxFetchAttempts-1 {
			break
		}

		log.Printf("%s: 429, retrying in %s", url, wait.Round(time.Second))
		time.Sleep(wait)
	}

	return nil, fmt.Errorf("%s: rate limited after %d attempts", url, maxFetchAttempts)
}

// tryFetch performs one GET. A non-zero wait means the request was rejected as
// rate limited and should be retried after that long; the body is nil then.
func tryFetch(c *http.Client, apiToken, url string, l *limiter) ([]byte, time.Duration, error) {
	res, err := request(c, apiToken, url)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	l.observe(res.Header)

	if res.StatusCode == http.StatusTooManyRequests {
		// A 429 without a usable Retry-After still has the reset headers to fall
		// back on, and failing that a wait long enough to clear a burst window.
		wait := max(l.burst.wait(), l.sustained.wait())

		if d, ok := retryAfter(res.Header.Get("Retry-After")); ok {
			wait = d
		}

		return nil, max(wait, time.Second), nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)

	return b, 0, err
}

// retryAfter parses either RFC 9110 form of Retry-After.
func retryAfter(value string) (time.Duration, bool) {
	if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second, true
	}

	if t, err := http.ParseTime(value); err == nil {
		return max(0, time.Until(t)), true
	}

	return 0, false
}

func request(c *http.Client, apiToken, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

	return c.Do(req)
}
