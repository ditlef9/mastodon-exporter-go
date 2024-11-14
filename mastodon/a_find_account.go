// mastodon/a_find_account.go
package mastodon

import (
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// SearchResponse is the structure for the Mastodon search API response
type SearchResponse struct {
	Accounts []models.Account `json:"accounts"`
}

// FindAccount takes a username and returns the account ID
func FindAccount(username string) (string, error) {
	// URL of the Mastodon API endpoint
	url := fmt.Sprintf("https://mastodon.social/api/v2/search?q=%s&type=accounts", username)

	// Send the GET request to the Mastodon API
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Unmarshal the JSON response
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", err
	}

	// Extract the account ID
	if len(searchResp.Accounts) > 0 {
		return searchResp.Accounts[0].ID, nil
	}
	return "", fmt.Errorf("a_find_account.FindAccount()::no account found for '%s'", username)
}
