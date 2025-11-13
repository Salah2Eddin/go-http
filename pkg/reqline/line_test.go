package reqline

import (
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/uri"
)

func TestNewRequestLine(t *testing.T) {
	u := uri.NewUri("/path")
	rl := NewRequestLine("GET", u, "HTTP/1.1")

	if rl.Method != "GET" {
		t.Errorf("Method = %q, want %q", rl.Method, "GET")
	}
	if rl.Uri != u {
		t.Errorf("Uri = %v, want %v", rl.Uri, u)
	}
	if rl.Version != "HTTP/1.1" {
		t.Errorf("Version = %q, want %q", rl.Version, "HTTP/1.1")
	}
}
