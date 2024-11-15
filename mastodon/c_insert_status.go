package mastodon

import (
	"database/sql"
	"ekeberg.com/mastodon-statuses-to-postgres-go/db"
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon/models"
	_ "github.com/lib/pq"
	"log"
)

// InsertStatus inserts a status into the messages_index table if it doesn't exist
// If the status exists and the content is different, it will update the content.
func InsertStatus(status *models.Status) (string, error) {
	log.Printf("c_insert_status.InsertStatus()::Inserting status with ID %s\n", status.ID)

	// Check if the status already exists using the external ID
	query := "SELECT msg_id, msg_content FROM messages_index WHERE msg_platform=$1 AND msg_external_id=$2"
	var sqMsgId string
	var dbContent string
	err := db.DB.QueryRow(query, "Mastodon", status.ID).Scan(&sqMsgId, &dbContent)

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("c_insert_status.InsertStatus()::Error checking if status exists: %v", err)
	}

	// If the status doesn't exist, insert the new status
	if sqMsgId == "" {
		insertQuery := `
		INSERT INTO messages_index (msg_platform, msg_external_id, msg_created_at, msg_language, msg_url, 
		                             msg_content, msg_external_account_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err := db.DB.Exec(insertQuery, "Mastodon", status.ID, status.CreatedAt, status.Language, status.URL, status.Content, status.AccountId)
		if err != nil {
			log.Fatalf("c_insert_status.InsertStatus()::Error inserting status: %v", err)
		} else {
			log.Printf("c_insert_status.InsertStatus()::Status with ID %s inserted successfully", status.ID)
		}

		// After inserting, get the msg_id
		err = db.DB.QueryRow("SELECT msg_id FROM messages_index WHERE msg_platform=$1 AND msg_external_id=$2", "Mastodon", status.ID).Scan(&sqMsgId)
		if err != nil {
			log.Fatalf("c_insert_status.InsertStatus()::Error retrieving msg_id for status: %v", err)
		}
	} else {
		// If the status exists, check if the content is different
		if dbContent != status.Content {
			log.Printf("c_insert_status.InsertStatus()::Status content has changed, updating msg_content for msg_id %s\n", sqMsgId)

			// Update the msg_content if it differs
			updateQuery := `
			UPDATE messages_index 
			SET msg_content=$1 
			WHERE msg_id=$2
			`
			_, err := db.DB.Exec(updateQuery, status.Content, sqMsgId)
			if err != nil {
				log.Fatalf("c_insert_status.InsertStatus()::Error updating status content: %v", err)
			} else {
				log.Printf("c_insert_status.InsertStatus()::Status content for msg_id %s updated successfully", sqMsgId)
			}
		}
	}

	// Return the msg_id of the inserted or existing status
	return sqMsgId, nil
}
