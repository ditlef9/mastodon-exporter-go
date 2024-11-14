package test

import (
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStatuses(t *testing.T) {
	// Create a mock HTTP client
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock a response from the Mastodon API
	httpmock.RegisterResponder("GET", "https://mastodon.social/api/v1/accounts/12345/statuses",
		httpmock.NewStringResponder(200, `[{"id": "123", "content": "Test status"}]`))

	// Call the function that makes the HTTP request
	statuses, err := mastodon.GetStatuses("12345")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that we got a status in the response
	assert.Len(t, statuses, 1)
	assert.Equal(t, statuses[0].ID, "123")
	assert.Equal(t, statuses[0].Content, "Test status")
}
