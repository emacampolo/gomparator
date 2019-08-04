package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-time.After(50 * time.Millisecond)
		}),
	)
	defer server.Close()
	c := New(Timeout(10 * time.Millisecond))
	c.client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		return false, nil
	}
	_, err := c.Fetch(server.URL, nil)

	assert.EqualError(t, err, fmt.Sprintf("Get %s: net/http: timeout awaiting response headers", server.URL))
}

func TestRetries(t *testing.T) {
	t.Parallel()
	var count int
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if count == 3 {
				w.WriteHeader(http.StatusOK)
				return
			}
			<-time.After(50 * time.Millisecond)
			count++
		}),
	)
	defer server.Close()
	c := New(Timeout(10 * time.Millisecond))
	res, _ := c.Fetch(server.URL, nil)
	assert.Equal(t, 200, res.StatusCode)
}

func TestRetryTimeout(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-time.After(50 * time.Millisecond)
		}),
	)
	defer server.Close()
	c := New(Timeout(10 * time.Millisecond))
	_, err := c.Fetch(server.URL, nil)

	assert.EqualError(t, err, fmt.Sprintf("GET %s giving up after 5 attempts", server.URL))
}

func TestResponseBodyCapture(t *testing.T) {
	t.Parallel()
	want := []byte("gomparator")
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(want)
		}),
	)
	defer server.Close()

	c := New()
	res, _ := c.Fetch(server.URL, nil)

	assert.Equal(t, want, res.Body)
}

func TestStatusCode(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer server.Close()
	c := New()
	res, _ := c.Fetch(server.URL, nil)

	assert.Equal(t, 400, res.StatusCode)
}
