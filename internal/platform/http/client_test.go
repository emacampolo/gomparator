package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	want := "net/http: timeout awaiting response headers"
	if got := err.Error(); !strings.HasSuffix(got, want) {
		t.Fatalf("want: '%v' in '%v'", want, got)
	}
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
	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestRetryTimeout(t *testing.T) {
	t.Parallel()
	var got int
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got++
			<-time.After(50 * time.Millisecond)
		}),
	)
	defer server.Close()
	c := New(Timeout(10 * time.Millisecond))
	_, err := c.Fetch(server.URL, nil)

	if want := 4; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	want := "giving up after 4 attempts"
	if got := err.Error(); !strings.HasSuffix(got, want) {
		t.Fatalf("want: '%v' in '%v'", want, got)
	}

}

func TestConnections(t *testing.T) {
	t.Parallel()
	c := New(Connections(23))
	got := c.client.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost
	if want := 23; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
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
	if got := res.Body; !bytes.Equal(got, want) {
		t.Fatalf("got: %v, want: %v", got, want)
	}
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
	if got, want := res.StatusCode, 400; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}
