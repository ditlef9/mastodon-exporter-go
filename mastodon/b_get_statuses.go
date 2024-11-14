package mastodon

import (
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetStatuses fetches the statuses for a given account ID
func GetStatuses(accountID string) ([]models.Status, error) {
	// Debug: Print the account ID being searched
	fmt.Printf("b_get_statuses.GetStatuses()::Searching statuses for account ID: %s\n", accountID)

	// URL of the Mastodon API endpoint
	url := fmt.Sprintf("https://mastodon.social/api/v1/accounts/%s/statuses", accountID)

	// Send the GET request to fetch statuses from the account
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("b_get_statuses.GetStatuses()::failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response code is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("b_get_statuses.GetStatuses()::failed to fetch statuses: %s", resp.Status)
	}

	// Read the entire response body for debugging
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("b_get_statuses.GetStatuses()::failed to read response body: %w", err)
	}

	// Print raw JSON response for debugging
	fmt.Println("b_get_statuses.GetStatuses()::Raw API Response:")
	fmt.Println("b_get_statuses.GetStatuses()::", string(body))

	// Decode the response JSON into a slice of Status objects
	var statuses []models.Status
	if err := json.Unmarshal(body, &statuses); err != nil {
		return nil, fmt.Errorf("b_get_statuses.GetStatuses()::failed to decode response body: %w", err)
	}

	// Inject accountID into each status
	for i := range statuses {
		statuses[i].AccountId = accountID
	}

	// Return the statuses
	return statuses, nil
}
