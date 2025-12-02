package timelib_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twistingmercury/docker-ci-build/timelib"
)

func TestNow(t *testing.T) {
	result := timelib.Now()
	assert.NotEmpty(t, result)
	// Verify it's a recent time (within 1 second)
	// Parse using local timezone to match time.Now().String() output
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", result[:19], time.Local)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), parsed, time.Second)
}
