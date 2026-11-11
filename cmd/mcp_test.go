package main

import (
	"testing"

	"github.com/pulseaiclub/phi/internal/mcp"
	"github.com/stretchr/testify/assert"
)

func TestFormatServerLine(t *testing.T) {
	cases := []struct {
		name   string
		server string
		cfg    mcp.ServerConfig
		want   string
	}{
		{
			name:   "stdio shows command and args",
			server: "fetch",
			cfg: mcp.ServerConfig{
				Transport: "stdio",
				Command:   []string{"npx", "-y"},
				Args:      []string{"@modelcontextprotocol/server-fetch"},
			},
			want: "fetch	npx -y @modelcontextprotocol/server-fetch",
		},
		{
			name:   "stdio default transport shows command and args",
			server: "fetch",
			cfg: mcp.ServerConfig{
				Command: []string{"node", "server.js"},
			},
			want: "fetch	node server.js",
		},
		{
			name:   "http shows configured url",
			server: "remote",
			cfg: mcp.ServerConfig{
				Transport: "http",
				URL:       "http://localhost:3000/mcp",
			},
			want: "remote	http://localhost:3000/mcp",
		},
		{
			name:   "streamable-http shows configured url",
			server: "stream",
			cfg: mcp.ServerConfig{
				Transport: "streamable-http",
				URL:       "https://api.example.com/mcp",
			},
			want: "stream	https://api.example.com/mcp",
		},
		{
			name:   "sse shows configured url",
			server: "events",
			cfg: mcp.ServerConfig{
				Transport: "sse",
				URL:       "https://api.example.com/sse",
			},
			want: "events	https://api.example.com/sse",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, formatServerLine(tc.server, tc.cfg))
		})
	}
}

func TestFormatServerListAlphabetical(t *testing.T) {
	servers := map[string]mcp.ServerConfig{
		"zeta": {
			Transport: "http",
			URL:       "https://zeta.example.com/mcp",
		},
		"alpha": {
			Command: []string{"node", "alpha.js"},
		},
		"middle": {
			Transport: "sse",
			URL:       "https://middle.example.com/sse",
		},
	}

	got := formatServerList(servers)
	want := "alpha\tnode alpha.js\n" +
		"middle\thttps://middle.example.com/sse\n" +
		"zeta\thttps://zeta.example.com/mcp\n"
	assert.Equal(t, want, got)
}

func TestFormatServerListOmitsAuthHeaders(t *testing.T) {
	servers := map[string]mcp.ServerConfig{
		"secure": {
			Transport: "http",
			URL:       "https://api.example.com/mcp",
			Headers: map[string]string{
				"Authorization": "Bearer secret-token",
			},
		},
	}

	got := formatServerList(servers)
	assert.Equal(t, "secure\thttps://api.example.com/mcp\n", got)
	assert.NotContains(t, got, "secret-token")
	assert.NotContains(t, got, "Authorization")
}
