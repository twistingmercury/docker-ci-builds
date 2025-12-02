package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAPIURL() string {
	if url := os.Getenv("API_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

func TestGetTime(t *testing.T) {
	// Make GET request to /api/time
	resp, err := http.Get(getAPIURL() + "/api/time")
	require.NoError(t, err, "failed to make request")
	defer resp.Body.Close()

	// Verify 200 status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response has "time" field
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "failed to decode response")

	assert.Contains(t, result, "time", "response should contain 'time' field")
}
