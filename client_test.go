package metron_test

import (
	"context"
	"embed"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/*
var fs embed.FS

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type requestMock struct {
	expectedURL         string
	expectedMethod      string
	responseStatus      int
	responseHeader      http.Header
	responseBody        string
	responseBodyFixture string
	validateBody        func(t *testing.T, body []byte)
	// expectedHeaders are asserted against the request. An empty value asserts
	// the header is absent, since Header.Get reports both the same way.
	expectedHeaders map[string]string
	// captureURL receives the request URL, for tests that build the expectation
	// rather than declaring it.
	captureURL *string
}

type testCase[T any] struct {
	id       int
	expected T
}

func testList[T, F any](
	t *testing.T,
	kind string,
	method func(*metron.Client, context.Context, *F, ...metron.RequestOption) iter.Seq2[T, error],
	testCases []testCase[T],
) {
	c := newTestClient(t, []requestMock{
		{expectedURL: fmt.Sprintf("https://metron.cloud/api/%s/?page=1", kind), responseBodyFixture: fmt.Sprintf("fixtures/%s_list_1.json", kind)},
		{expectedURL: fmt.Sprintf("https://metron.cloud/api/%s/?page=2", kind), responseBodyFixture: fmt.Sprintf("fixtures/%s_list_2.json", kind)},
	})

	resources := make([]T, 0, 4)

	for res, err := range method(c, context.Background(), nil) {
		require.NoError(t, err)
		resources = append(resources, res)
	}

	require.Len(t, resources, len(testCases))

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%T #%d", tc.expected, tc.id), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, resources[i])
		})
	}
}

func testListByID[T any](
	t *testing.T,
	kind string,
	id int,
	listKind string,
	method func(*metron.Client, context.Context, int, ...metron.ConditionalOption) iter.Seq2[T, error],
	testCases []testCase[T],
) {
	c := newTestClient(t, []requestMock{
		{
			expectedURL:         fmt.Sprintf("https://metron.cloud/api/%s/%d/%s_list/?page=1", kind, id, listKind),
			responseBodyFixture: fmt.Sprintf("fixtures/%s_%d_%s_list_1.json", kind, id, listKind),
		},
		{
			expectedURL:         fmt.Sprintf("https://metron.cloud/api/%s/%d/%s_list/?page=2", kind, id, listKind),
			responseBodyFixture: fmt.Sprintf("fixtures/%s_%d_%s_list_2.json", kind, id, listKind),
		},
	})

	resources := make([]T, 0, 4)

	for res, err := range method(c, context.Background(), id) {
		require.NoError(t, err)
		resources = append(resources, res)
	}

	require.Len(t, resources, len(testCases))

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%T #%d", tc.expected, tc.id), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, resources[i])
		})
	}
}

func testByID[T any](
	t *testing.T,
	kind string,
	method func(*metron.Client, context.Context, int, ...metron.ConditionalOption) (T, error),
	testCases []testCase[T],
) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T #%d", tc.expected, tc.id), func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, []requestMock{
				{
					expectedURL:         fmt.Sprintf("https://metron.cloud/api/%s/%d/", kind, tc.id),
					responseBodyFixture: fmt.Sprintf("fixtures/%s_%d.json", kind, tc.id),
				},
			})

			v, err := method(c, context.Background(), tc.id)

			require.NoError(t, err)
			require.Equal(t, tc.expected, v)
		})
	}
}

func newTestClient(t *testing.T, mocks []requestMock) *metron.Client {
	c, err := metron.NewClient("foobar", metron.WithClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			m := mocks[0]
			mocks = mocks[1:]

			if m.captureURL != nil {
				*m.captureURL = req.URL.String()
			}

			if m.expectedURL != "" {
				require.Equal(t, m.expectedURL, req.URL.String())
			}

			if m.expectedMethod != "" {
				require.Equal(t, m.expectedMethod, req.Method)
			}

			require.Equal(t, "Bearer foobar", req.Header.Get("Authorization"))

			for name, value := range m.expectedHeaders {
				require.Equal(t, value, req.Header.Get(name), name)
			}

			if m.validateBody != nil {
				bodyBytes, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				m.validateBody(t, bodyBytes)
			}

			var body io.ReadCloser

			if m.responseBodyFixture != "" {
				f, err := fs.Open(m.responseBodyFixture)
				require.NoError(t, err)

				body = f
			} else {
				body = io.NopCloser(strings.NewReader(m.responseBody))
			}

			status := m.responseStatus
			if status == 0 {
				status = http.StatusOK
			}

			header := m.responseHeader
			if header == nil {
				header = make(http.Header)
			}

			return &http.Response{
				StatusCode: status,
				Body:       body,
				Header:     header,
			}
		}),
	}))
	require.NoError(t, err)

	return c
}

func parseDate(t *testing.T, dateString string) civil.Date {
	v, err := civil.ParseDate(dateString)

	require.NoError(t, err)

	return v
}

func parseTime(t *testing.T, timeString string) time.Time {
	v, err := time.Parse(time.RFC3339Nano, timeString)

	require.NoError(t, err)

	return v
}

func parseURL(t *testing.T, urlString string) url.URL {
	v, err := url.Parse(urlString)

	require.NoError(t, err)

	return *v
}

func TestWithBaseURL(t *testing.T) {
	t.Parallel()

	var got string

	c, err := metron.NewClient("foobar",
		metron.WithBaseURL("http://localhost:8080/metron"),
		metron.WithClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) *http.Response {
				got = req.URL.String()

				f, openErr := fs.Open("fixtures/arc_659.json")
				require.NoError(t, openErr)

				return &http.Response{StatusCode: http.StatusOK, Body: f, Header: make(http.Header)}
			}),
		}),
	)
	require.NoError(t, err)

	_, err = c.ArcByID(context.Background(), 659)
	require.NoError(t, err)

	require.Equal(t, "http://localhost:8080/metron/api/arc/659/", got)
}

// TestListIteratorIsSafeForConcurrentIteration pins the reason paginate builds
// params per page: the sequence is a value, not a cursor.
func TestListIteratorIsSafeForConcurrentIteration(t *testing.T) {
	t.Parallel()

	c, err := metron.NewClient("foobar", metron.WithClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			body := `{"count":0,"next":null,"previous":null,"results":[]}`
			if req.URL.Query().Get("page") == "1" {
				body = `{"count":0,"next":"https://metron.cloud/api/issue/?page=2","previous":null,"results":[]}`
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		}),
	}))
	require.NoError(t, err)

	seq := c.Issues(context.Background(), &metron.IssueFilters{PublisherName: "marvel"})

	var wg sync.WaitGroup

	errs := make(chan error, 16)

	for range 2 {
		wg.Go(func() {
			for _, iterErr := range seq {
				if iterErr != nil {
					errs <- iterErr
				}
			}
		})
	}

	wg.Wait()
	close(errs)

	for iterErr := range errs {
		require.NoError(t, iterErr)
	}
}

func TestWithBaseURLInvalid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		baseURL       string
		expectedError string
	}{
		{
			name:          "javascript scheme",
			baseURL:       "javascript:alert(1)",
			expectedError: "scheme must be http or https",
		},
		{
			name:          "ftp scheme",
			baseURL:       "ftp://metron.cloud",
			expectedError: "scheme must be http or https",
		},
		{
			name:          "relative path",
			baseURL:       "/api",
			expectedError: "scheme must be http or https",
		},
		{
			name:          "no host",
			baseURL:       "https://",
			expectedError: "missing host",
		},
		{
			name:          "empty",
			baseURL:       "",
			expectedError: "scheme must be http or https",
		},
		{
			name:          "unparseable",
			baseURL:       "https://metron.cloud/%zz",
			expectedError: "invalid URL escape",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := metron.NewClient("foobar", metron.WithBaseURL(tc.baseURL))

			require.Nil(t, c)
			require.ErrorContains(t, err, tc.expectedError)
		})
	}
}

// TestUserAgent checks the User-Agent built from build info. Deps is empty in
// our own test binary (golang/go#68045), so it falls back to "devel".
func TestUserAgent(t *testing.T) {
	t.Parallel()

	const self = "go-metron/devel"

	testCases := []struct {
		name     string
		options  []metron.Option
		expected string
	}{
		{
			name:     "default",
			expected: self,
		},
		{
			name:     "with prefix",
			options:  []metron.Option{metron.WithUserAgent("myapp/1.2.3")},
			expected: "myapp/1.2.3 " + self,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got string

			options := append([]metron.Option{
				metron.WithClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) *http.Response {
						got = req.Header.Get("User-Agent")

						f, err := fs.Open("fixtures/arc_659.json")
						require.NoError(t, err)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       f,
							Header:     make(http.Header),
						}
					}),
				}),
			}, tc.options...)

			c, err := metron.NewClient("foobar", options...)
			require.NoError(t, err)

			_, err = c.ArcByID(context.Background(), 659)
			require.NoError(t, err)

			require.Equal(t, tc.expected, got)
		})
	}
}
