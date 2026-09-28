package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate_UnmarshalJSON(t *testing.T) {
	type TestReq struct {
		ReleaseDate Date `json:"release_date"`
	}

	// Test YYYY-MM-DD
	var req1 TestReq
	err := json.Unmarshal([]byte(`{"release_date":"2026-09-25"}`), &req1)
	assert.NoError(t, err)
	assert.Equal(t, 2026, time.Time(req1.ReleaseDate).Year())
	assert.Equal(t, time.September, time.Time(req1.ReleaseDate).Month())
	assert.Equal(t, 25, time.Time(req1.ReleaseDate).Day())

	// Test RFC3339
	var req2 TestReq
	err = json.Unmarshal([]byte(`{"release_date":"2026-09-25T00:00:00Z"}`), &req2)
	assert.NoError(t, err)
	assert.Equal(t, 2026, time.Time(req2.ReleaseDate).Year())
}
