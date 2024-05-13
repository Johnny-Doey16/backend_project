package services

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/worker"
	// Import other necessary packages, such as for UUID handling and database operations.
)

func extractUsernames(content string) []string {
	mentionRegex := regexp.MustCompile(`@(\w+)`)
	matches := mentionRegex.FindAllStringSubmatch(content, -1)

	usernameSet := make(map[string]struct{})
	for _, match := range matches {
		username := strings.ToLower(match[1])
		usernameSet[username] = struct{}{}
	}

	usernames := make([]string, 0, len(usernameSet))
	for username := range usernameSet {
		usernames = append(usernames, username)
	}
	return usernames
}

func DetectMentions(ctx context.Context, td *worker.TaskDistributor, creator, content string, qtx *sqlc.Queries, tx *sql.Tx, authorId, postID uuid.UUID) error {
	// Extract unique usernames from the content.
	usernames := extractUsernames(content)

	// If no usernames, nothing to do.
	if len(usernames) == 0 {
		return nil
	}

	// Get user IDs from usernames.
	userIDs, err := GetUserIDsFromUsernames(tx, creator, usernames)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("1. %+v", err)
	}

	// // Check if the mention already exists.
	// If the mention doesn't exist, create it and notify the user.
	if err := createMention(ctx, qtx, postID, userIDs); err != nil {
		tx.Rollback()
		return fmt.Errorf("2. %+v", err)
	}
	if err := notifyUserOfMention(ctx, td, creator, userIDs, authorId, postID); err != nil {
		tx.Rollback()
		return fmt.Errorf("3. %+v", err)
	}

	// Commit the transaction.
	return nil //tx.Commit()
}

// GetUserIDsFromUsernames queries the database to get user IDs for the given usernames.
func GetUserIDsFromUsernames(tx *sql.Tx, creator string, usernames []string) ([]uuid.UUID, error) {
	if len(usernames) == 0 {
		return []uuid.UUID{}, nil
	}

	// Create a set of unique usernames to eliminate duplicates.
	usernameSet := make(map[string]struct{})
	for _, username := range usernames {
		if strings.EqualFold(username, creator) {
			continue
		}
		usernameSet[strings.ToLower(username)] = struct{}{} //
		// usernameSet[username] = struct{}{}
	}

	// Convert the set back to a slice for use in the query.
	uniqueUsernames := make([]string, 0, len(usernameSet))
	for username := range usernameSet {
		uniqueUsernames = append(uniqueUsernames, username)
	}

	// Prepare the SQL query to retrieve user IDs.
	// NOTE: This uses PostgreSQL syntax; adapt the placeholder format if using another database.
	query := `SELECT id FROM authentications WHERE LOWER(username) = ANY($1)`
	rows, err := tx.Query(query, pq.Array(uniqueUsernames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	// Check for any error that occurred during iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil

	// return qtx.GetUidsFromUsername(ctx, uniqueUsernames)
}

// mentionExists checks if a mention already exists in the database.
// ! Use in edit post
func mentionExists(ctx context.Context, tx *sql.Tx, postID uuid.UUID, userID uuid.UUID) (bool, error) {
	// Prepare the SQL query to check if the mention exists.
	query := `SELECT EXISTS(SELECT 1 FROM post_mentions WHERE post_id = $1 AND mentioned_user_id = $2)`
	var exists bool
	err := tx.QueryRowContext(ctx, query, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// createMention inserts a mention into the post_mentions table within the transaction.
func createMention(ctx context.Context, qtx *sqlc.Queries, postID uuid.UUID, userID []uuid.UUID) error {
	// Prepare the SQL statement for inserting a new mention.
	// query := `INSERT INTO post_mentions (post_id, mentioned_user_id) VALUES ($1, $2)`

	// // Execute the SQL statement using the provided transaction.
	// _, err := tx.ExecContext(ctx, query, postID, userID)
	// if err != nil {
	// 	return err
	// }

	// return nil

	return qtx.CreatePostMention(ctx, sqlc.CreatePostMentionParams{
		PostID:  postID,
		Column1: userID,
	})
}

// sends a notification to the user about being mentioned in a post.
func notifyUserOfMention(ctx context.Context, td *worker.TaskDistributor, creator string, userIDs []uuid.UUID, authorId, postID uuid.UUID) error {

	if len(userIDs) > 0 {
		notificationMessage := fmt.Sprintf("%s mentioned you in their post", creator)
		SendNotification(*td, ctx, authorId, postID, userIDs, constants.NotificationPostMention, "Mention", "", notificationMessage, "")
	}

	return nil
}
