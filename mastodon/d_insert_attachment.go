package mastodon

import (
	"database/sql"
	"ekeberg.com/mastodon-statuses-to-postgres-go/db"
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon/models"
	_ "github.com/lib/pq"
	"log"
)

// InsertAttachment inserts a media attachment into the messages_attachments table if it doesn't exist
// If the media exists and fields have changed, it will update the fields.
func InsertAttachment(media *models.MediaAttachment, msgID string) {
	log.Printf("d_insert_attachment.InsertAttachment()::Inserting media with ID %s\n", media.ID)

	// Check if the media attachment already exists using the attachment_external_id
	query := "SELECT attachment_id, attachment_url, attachment_type, attachment_meta_description FROM messages_attachments WHERE attachment_external_id=$1"
	var sqlAttachmentID, sqlURL, sqlType, sqlDescription string
	err := db.DB.QueryRow(query, media.ID).Scan(&sqlAttachmentID, &sqlURL, &sqlType, &sqlDescription)

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("d_insert_attachment.InsertAttachment()::Error checking if media exists: %v", err)
	}

	// If the media doesn't exist, insert it
	if sqlAttachmentID == "" {
		insertQuery := `
		INSERT INTO messages_attachments (attachment_msg_id, attachment_external_id, attachment_url, attachment_type, attachment_meta_description)
		VALUES ($1, $2, $3, $4, $5)
		`
		_, err := db.DB.Exec(insertQuery, msgID, media.ID, media.URL, media.Type, media.MetaDescription)
		if err != nil {
			log.Fatalf("d_insert_attachment.InsertAttachment()::Error inserting media: %v", err)
		} else {
			log.Printf("d_insert_attachment.InsertAttachment()::Media with ID %s inserted successfully", media.ID)
		}
	} else {
		// If the media exists, check if any field is different
		updateNeeded := false
		if sqlURL != media.URL {
			log.Printf("d_insert_attachment.InsertAttachment()::Media URL has changed for media ID %s", media.ID)
			updateNeeded = true
		}
		if sqlType != media.Type {
			log.Printf("d_insert_attachment.InsertAttachment()::Media type has changed for media ID %s", media.ID)
			updateNeeded = true
		}
		if sqlDescription != media.MetaDescription {
			log.Printf("d_insert_attachment.InsertAttachment()::Media description has changed for media ID %s", media.ID)
			updateNeeded = true
		}

		// If any field is different, update the media attachment
		if updateNeeded {
			updateQuery := `
			UPDATE messages_attachments
			SET attachment_url=$1, attachment_type=$2, attachment_meta_description=$3
			WHERE attachment_external_id=$4
			`
			_, err := db.DB.Exec(updateQuery, media.URL, media.Type, media.MetaDescription, media.ID)
			if err != nil {
				log.Fatalf("d_insert_attachment.InsertAttachment()::Error updating media: %v", err)
			} else {
				log.Printf("d_insert_attachment.InsertAttachment()::Media with ID %s updated successfully", media.ID)
			}
		} else {
			log.Printf("d_insert_attachment.InsertAttachment()::No changes detected for media with ID %s, skipping update", media.ID)
		}
	}
}
