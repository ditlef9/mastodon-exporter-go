// mastodon/models/status.go

package models

// Define the struct to match the JSON response for statuses
type Status struct {
	ID               string            `json:"id"`
	CreatedAt        string            `json:"created_at"`
	Language         string            `json:"language"`
	URL              string            `json:"url"`
	Content          string            `json:"content"`
	AccountId        string            `json:"account.id"`
	MediaAttachments []MediaAttachment `json:"media_attachments"`
}
