// main.go
package main

import (
	"ekeberg.com/mastodon-statuses-to-postgres-go/db"
	"ekeberg.com/mastodon-statuses-to-postgres-go/mastodon"
	"fmt"
	"log"
	"os"
)

// Main
func main() {
	// SQLite connection
	db.InitDB()

	// Settings
	mastodonUsername := os.Getenv("MASTODON_USERNAME") // User's Mastodon handle (without @domain)

	// Get account id
	mastodonAccountID, err := mastodon.FindAccount(mastodonUsername)
	if err != nil {
		log.Fatal(err)
	}

	// Get statuses
	mastodonStatuses, err := mastodon.GetStatuses(mastodonAccountID)
	if err != nil {
		log.Fatal(err)
	}

	// Print statuses with additional information
	for _, status := range mastodonStatuses {
		// Print
		fmt.Printf("main.main()::Status ID: %s - %s\n", status.ID, status.Content)

		// Insert status into the database if it doesn't exist
		msgID, err := mastodon.InsertStatus(&status)
		if err != nil {
			log.Fatalf("main.main()::Error inserting status: %v", err)
		}
		fmt.Printf("main.main()::Inserted/Existing Status with msg_id: %s\n", msgID)

		// Debug:
		//
		// fmt.Printf("main.main()::Created at: %s\n", status.CreatedAt)
		// fmt.Printf("main.main()::Language: %s\n", status.Language)
		// fmt.Printf("main.main()::URL: %s\n", status.URL)
		// fmt.Printf("main.main()::Content: %s\n", status.Content)
		// fmt.Printf("main.main()::Account ID: %s\n", status.AccountId)

		// Print media attachments
		if len(status.MediaAttachments) > 0 {
			for _, media := range status.MediaAttachments {
				// Insert into database
				mastodon.InsertAttachment(&media, msgID)

				// Debug:
				// fmt.Printf("  Media ID: %s\n", media.ID)
				// fmt.Printf("  Media Type: %s\n", media.Type)
				// fmt.Printf("  Media URL: %s\n", media.URL)
				// fmt.Printf("  Description: %s\n", media.MetaDescription)
			}
		}

		// Add a separator between statuses
		fmt.Println("--------------------------")
	}

	// Clean up old statuses and attachments
	mastodon.Cleanup(mastodonStatuses)
}
