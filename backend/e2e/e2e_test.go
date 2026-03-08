//go:build e2e

package e2e

import (
	"os"
	"strings"
)

var baseURL = "http://localhost:3000"

func init() {
	if u := os.Getenv("E2E_BASE_URL"); u != "" {
		baseURL = strings.TrimRight(u, "/")
	}
}
