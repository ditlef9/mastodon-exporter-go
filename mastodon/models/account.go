// mastodon/models/account.go

package models

// Define the struct to match the JSON response
type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
}
