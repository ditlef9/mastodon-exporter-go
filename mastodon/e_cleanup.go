package mastodon

import (
	"ekeberg.com/mastodon-statuses-to-postgres-go/db"
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon/models"
	"log"
)

// Cleanup function to remove old statuses and their media attachments
func Cleanup(mastodonStatuses []models.Status) {
	log.Println("e_cleanup.Cleanup()::Cleaning up old statuses and attachments...")

	// Step 1: Fetch the last 5 entries from the database
	query := `SELECT msg_id, msg_external_id FROM messages_index ORDER BY msg_id DESC LIMIT 5`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("e_cleanup.Cleanup()::Error fetching last 5 statuses: %v", err)
	}
	defer rows.Close()

	// Store the msg_external_id of the last 5 statuses
	var dbStatuses []string
	for rows.Next() {
		var msgID, msgExternalID string
		if err := rows.Scan(&msgID, &msgExternalID); err != nil {
			log.Fatalf("e_cleanup.Cleanup()::Error scanning row: %v", err)
		}
		dbStatuses = append(dbStatuses, msgExternalID)
	}

	// Step 2: Compare fetched statuses with mastodonStatuses
	for _, status := range mastodonStatuses {
		// Remove the status from dbStatuses if it exists in mastodonStatuses
		for i, dbStatus := range dbStatuses {
			if dbStatus == status.ID {
				// Found a match, so we remove it from the list to avoid deletion
				dbStatuses = append(dbStatuses[:i], dbStatuses[i+1:]...)
				break
			}
		}
	}

	// Step 3: Delete statuses that are not in mastodonStatuses
	for _, dbStatus := range dbStatuses {
		log.Printf("e_cleanup.Cleanup()::Deleting status with msg_external_id: %s\n", dbStatus)

		// Step 3.1: Delete the associated media attachments
		deleteAttachmentsQuery := `DELETE FROM messages_attachments WHERE attachment_msg_id IN (SELECT msg_id FROM messages_index WHERE msg_external_id = $1)`
		_, err := db.DB.Exec(deleteAttachmentsQuery, dbStatus)
		if err != nil {
			log.Printf("e_cleanup.Cleanup()::Error deleting attachments for status %s: %v", dbStatus, err)
		} else {
			log.Printf("e_cleanup.Cleanup()::Attachments for status %s deleted successfully", dbStatus)
		}

		// Step 3.2: Delete the status from the messages_index table
		deleteStatusQuery := `DELETE FROM messages_index WHERE msg_external_id = $1`
		_, err = db.DB.Exec(deleteStatusQuery, dbStatus)
		if err != nil {
			log.Printf("e_cleanup.Cleanup()::Error deleting status %s: %v", dbStatus, err)
		} else {
			log.Printf("e_cleanup.Cleanup()::Status with msg_external_id %s deleted successfully", dbStatus)
		}
	}
}
