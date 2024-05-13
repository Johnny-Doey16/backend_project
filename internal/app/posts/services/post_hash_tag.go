package services

import (
	"context"
	"database/sql"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// extractHashtags finds all unique hashtags in the given content.
// TODO: Check string case
func extractHashtags(content string) []string {
	hashtagRegex := regexp.MustCompile(`#(\w+)`)
	matches := hashtagRegex.FindAllStringSubmatch(content, -1)

	hashtagSet := make(map[string]struct{})
	for _, match := range matches {
		hashtag := strings.ToLower(match[1])
		hashtagSet[hashtag] = struct{}{}
	}

	hashtags := make([]string, 0, len(hashtagSet))
	for hashtag := range hashtagSet {
		hashtags = append(hashtags, hashtag)
	}
	return hashtags
}

// createHashtag inserts a new hashtag into the hashtag table if it doesn't exist.
func createHashtag(ctx context.Context, tx *sql.Tx, hashtag string) (int, error) {
	// Convert hashtag to lowercase
	hashtagLower := strings.ToLower(hashtag)

	var id int
	err := tx.QueryRowContext(ctx, `INSERT INTO hashtag (hash_tag) VALUES ($1) ON CONFLICT (hash_tag) DO UPDATE SET hash_tag=EXCLUDED.hash_tag RETURNING id`, hashtagLower).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// getHashtagIDsFromHashtags converts hashtag strings to their IDs.
func getHashtagIDsFromHashtags(ctx context.Context, tx *sql.Tx, hashtags []string) ([]int, error) {
	ids := make([]int, 0, len(hashtags))
	for _, hashtag := range hashtags {
		id, err := createHashtag(ctx, tx, hashtag)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// createPostHashtag associates a post with a hashtag in the post_hashtag table.
func createPostHashtag(ctx context.Context, tx *sql.Tx, postID uuid.UUID, hashtagID int) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO post_hashtag (post_id, hashtag_id) VALUES ($1, $2)`, postID, hashtagID)
	if err != nil {
		return err
	}
	return nil
}

// ProcessHashtags is the main function that extracts hashtags from content and creates the necessary associations.
func ProcessHashtags(ctx context.Context, content string, tx *sql.Tx, postID uuid.UUID) error {
	hashtags := extractHashtags(content)
	hashtagIDs, err := getHashtagIDsFromHashtags(ctx, tx, hashtags)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, hashtagID := range hashtagIDs {
		if err := createPostHashtag(ctx, tx, postID, hashtagID); err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}
