package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getApiURL() string {
	if url := os.Getenv("API_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

func TestGetUUID(t *testing.T) {
	resp, err := http.Get(getApiURL() + "/api/uuid")
	require.NoError(t, err, "failed to make request")
	defer func() {
		_ = resp.Body.Close()
	}()

	// Verify 200 status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response has "time" field
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "failed to decode response")
	assert.Contains(t, result, "uuid", "response should contain 'uuid' field")
	empty := uuid.UUID{}
	value := result["uuid"]
	t.Logf("%s\n", empty.String())
	t.Logf("%s\n", value)
	assert.NotEqual(t, empty, value)
}
