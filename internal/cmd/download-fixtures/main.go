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
	"time"
)

func main() {
	apiToken := flag.String("api-token", "", "api-token")

	flag.Parse()
	c := &http.Client{Timeout: time.Second * 10}

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
		if err := makeFixture(c, *apiToken, v.URL, v.FileName); err != nil {
			log.Println(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func makeFixture(c *http.Client, apiToken, url, fileName string) error {
	res, err := request(c, apiToken, url)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer

	if err = json.Indent(&prettyJSON, b, "", "  "); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("fixtures/", fileName), prettyJSON.Bytes(), 0o600)
}

func request(c *http.Client, apiToken, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

	return c.Do(req)
}
