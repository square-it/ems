package ems

import (
	"testing"
)

func assertNewClient(c *client, t *testing.T) {
	if c == nil {
		t.Fatalf("ops is nil")
	}

	if c.options.serverUrl.Host != "127.0.0.1:7222" {
		t.Fatalf("bad server host")
	}

	if c.options.serverUrl.Scheme != "tcp" {
		t.Fatalf("bad server scheme")
	}

	if c.options.username != "admin" {
		t.Fatalf("bad username")
	}

	if c.options.password != "" {
		t.Fatalf("bad password")
	}
}
