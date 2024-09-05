package metron

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type cacheable interface {
	modified() time.Time
}

type cacheEntry[T any] struct {
	Resource  T         `json:"resource"`
	ExpiresAt time.Time `json:"expires_at"`
}

func cache[T any](c *Client, r *http.Request) (T, error) {
	var v T
	var err error

	cacheDir := c.cacheDir

	if cacheDir == "" {
		cacheDir, err = defaultCacheDir()

		if err != nil {
			return v, err
		}
	}

	key := cacheKey(r, cacheDir)

	_, ok := interface{}(&v).(cacheable)

	if c.enableCaching && ok {
		if v, err = cacheGet[T](key); err != nil {
			if !os.IsNotExist(err) && !errors.Is(err, errCacheExpired) {
				return v, err
			}
		} else {
			return v, nil
		}
	}

	v, err = do[T](c, r)

	if err != nil {
		return v, err
	}

	if c.enableCaching && ok {
		if err = cachePut(key, v); err != nil {
			return v, err
		}
	}

	return v, err
}

var errCacheExpired = errors.New("cache entry expired")

func cacheGet[T any](cacheKey string) (T, error) {
	var v T

	f, err := os.Open(cacheKey)

	if err != nil {
		return v, err
	}

	var b []byte

	b, err = io.ReadAll(f)

	if err != nil {
		return v, err
	}

	var c cacheEntry[T]

	if err = json.Unmarshal(b, &c); err != nil {
		return v, err
	}

	if c.ExpiresAt.Before(time.Now()) {
		return v, errCacheExpired
	}

	return c.Resource, nil
}

func cacheKey(r *http.Request, cacheDir string) string {
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(r.URL.String())))

	return filepath.Join(cacheDir, sum[0:3], sum[3:6], sum[6:9], sum[9:12], sum[12:15], sum[15:18], sum)
}

func cachePut[T any](cacheKey string, v T) error {
	if err := os.MkdirAll(filepath.Dir(cacheKey), 0755); err != nil {
		return err
	}

	c, ok := interface{}(&v).(cacheable)

	if !ok {
		return errors.New("invalid cache entry")
	}

	maxAge := time.Now().Add(time.Duration(float64(time.Since(c.modified())) * 0.1))
	b, err := json.Marshal(cacheEntry[T]{
		Resource:  v,
		ExpiresAt: maxAge,
	})

	if err != nil {
		return err
	}

	if err = os.WriteFile(cacheKey, b, 0644); err != nil {
		return err
	}

	return nil
}

func defaultCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "metron"), nil
}
