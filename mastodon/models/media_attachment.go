// mastodon/models/media_attachment.go

package models

// MediaAttachment represents a media file attached to a status
type MediaAttachment struct {
	ID              string `json:"id"`
	Type            string `json:"type"` // Image, video, etc.
	URL             string `json:"url"`
	MetaDescription string `json:"description"`
}
