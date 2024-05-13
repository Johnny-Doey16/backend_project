package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	fbservice "github.com/steve-mir/diivix_backend/fb-service"
)

const (
	TaskSendNotification  = "task:send_notification"
	likeUniqueErr         = `pq: duplicate key value violates unique constraint "unique_like_per_user_per_post"`
	commentUniqueErr      = `pq: duplicate key value violates unique constraint "unique_comment_per_user_per_post"`
	churchAnnUniqueErr    = `pq: duplicate key value violates unique constraint "unique_user_per_announcement"`
	mentionUniqueErr      = `pq: duplicate key value violates unique constraint "unique_mention_per_user_per_post"`
	prayerInviteUniqueErr = `pq: duplicate key value violates unique constraint "unique_prayer_invite_per_user_per_room"`
)

type PayloadSendNotification struct {
	UserID   []uuid.UUID `json:"user_ids"`
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Body     string      `json:"body"`
	ImageUrl string      `json:"image_url"`
	PostID   uuid.UUID   `json:"post_id"`
	AuthorID uuid.UUID   `json:"author_id"`
	RoomID   string      `json:"room_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendNotification(
	ctx context.Context,
	payload *PayloadSendNotification,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendNotification, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued notification task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendNotification(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendNotification
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	// log.Debug().Msg("PAyload "+ string(payload.UserID[]))

	// Retrieve the user's notification details from the store using `payload.UserID`
	fcmTokens, err := processor.store.GetFCMTokenInSession(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("failed to get users' tokens: %w", err)
	}

	// ! Extract tokens from sql list
	tokens := make([]string, len(fcmTokens))
	for i, token := range fcmTokens {
		if token.String == "" {
			continue
		}
		tokens[i] = token.String
	}

	// Create transaction
	tx, err := processor.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := processor.store.WithTx(tx)

	notificationData := map[string]string{}

	// TODO: Refactor with Strategy pattern
	switch payload.Type {
	case constants.NotificationFollow:
		err = createFollowNotification(ctx, qtx, tx, payload)
		notificationData["activity"] = "FollowActivity"
		notificationData["author"] = payload.AuthorID.String()
	case constants.NotificationPostComment:
		err = createCommentNotification(ctx, qtx, tx, payload)
		notificationData["activity"] = "CommentActivity"
		notificationData["author"] = payload.AuthorID.String()
		notificationData["postId"] = payload.PostID.String()
	case constants.NotificationPostLike:
		err = createPostLikeNotification(ctx, qtx, tx, payload)
		notificationData["activity"] = "PostLikeActivity"
		notificationData["author"] = payload.AuthorID.String()
		notificationData["postId"] = payload.PostID.String()
	case constants.NotificationChurchAnnouncement:
		err = createAnnouncementNotification(ctx, qtx, tx, payload)
		notificationData["activity"] = "ChurchAnnouncementActivity"
		notificationData["author"] = payload.AuthorID.String()
	case constants.NotificationPostMention:
		err = createPostMentionNotification(ctx, qtx, tx, payload)
		notificationData["activity"] = "PostMentionActivity"
		notificationData["author"] = payload.AuthorID.String()
		notificationData["postId"] = payload.PostID.String()
	case constants.NotificationPrayerInvite:
		err = createPrayerInviteNotification(ctx, qtx, tx, payload)
		notificationData["activity"] = "PrayerActivity"
		notificationData["author"] = payload.AuthorID.String()
	}

	if err != nil {
		if err.Error() == likeUniqueErr || err.Error() == churchAnnUniqueErr || err.Error() == commentUniqueErr || err.Error() == mentionUniqueErr || err.Error() == prayerInviteUniqueErr {
			log.Error().Msg("Skipping error " + err.Error())
			return nil
		}
		return err
	}

	// Send the notification using the details retrieved
	err = fbservice.SendPushNotificationMulti(tokens, payload.Title, payload.Body, payload.ImageUrl, notificationData)
	if err != nil {
		return err
	}

	// Log the notification action
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("user_id", payload.UserID[0].String()).Msg("processed notification task")

	return nil
}

func createFollowNotification(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, payload PayloadSendNotification) error {
	id, err := qtx.CreateNotification(ctx, sqlc.CreateNotificationParams{
		Type:    payload.Type,
		Column1: payload.UserID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = qtx.CreateNotificationFollow(ctx, sqlc.CreateNotificationFollowParams{
		AuthorID:        payload.AuthorID,
		NotificationID:  id[0].ID,
		FollowingUserID: payload.UserID[0],
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createCommentNotification(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, payload PayloadSendNotification) error {
	id, err := qtx.CreateNotification(ctx, sqlc.CreateNotificationParams{
		Type:    payload.Type,
		Column1: payload.UserID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = qtx.CreateNotificationPostComment(ctx, sqlc.CreateNotificationPostCommentParams{
		AuthorID:       payload.AuthorID,
		NotificationID: id[0].ID,
		PostID:         payload.PostID,
		CommentUserID:  payload.UserID[0],
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createPostLikeNotification(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, payload PayloadSendNotification) error {
	id, err := qtx.CreateNotification(ctx, sqlc.CreateNotificationParams{
		Type:    payload.Type,
		Column1: payload.UserID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = qtx.CreateNotificationPostLike(ctx, sqlc.CreateNotificationPostLikeParams{
		AuthorID:       payload.AuthorID,
		NotificationID: id[0].ID,
		PostID:         payload.PostID,
		LikeUserID:     payload.UserID[0],
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createAnnouncementNotification(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, payload PayloadSendNotification) error {
	notifications, err := qtx.CreateNotification(ctx, sqlc.CreateNotificationParams{
		Type:    payload.Type,
		Column1: payload.UserID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	// Extract notification IDs into a slice of int32
	ids := make([]int32, len(notifications))
	uids := make([]uuid.UUID, len(notifications)) // TODO: consider removing to just use payload.UserId
	for i, notification := range notifications {
		ids[i] = notification.ID
		uids[i] = notification.UserID
	}

	// Insert notification announcements efficiently
	err = qtx.CreateNotificationAnnouncement(ctx, sqlc.CreateNotificationAnnouncementParams{
		AuthorID: payload.AuthorID,
		NewsID:   payload.PostID,
		Column1:  uids,
		Column2:  ids,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createPostMentionNotification(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, payload PayloadSendNotification) error {
	notifications, err := qtx.CreateNotification(ctx, sqlc.CreateNotificationParams{
		Type:    payload.Type,
		Column1: payload.UserID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	// Extract notification IDs into a slice of int32
	ids := make([]int32, len(notifications))
	uids := make([]uuid.UUID, len(notifications)) // TODO: consider removing to just use payload.UserId
	for i, notification := range notifications {
		ids[i] = notification.ID
		uids[i] = notification.UserID
	}

	err = qtx.CreateNotificationPostMention(ctx, sqlc.CreateNotificationPostMentionParams{
		AuthorID: payload.AuthorID,
		PostID:   payload.PostID,
		Column1:  uids,
		Column2:  ids,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createPrayerInviteNotification(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, payload PayloadSendNotification) error {
	notifications, err := qtx.CreateNotification(ctx, sqlc.CreateNotificationParams{
		Type:    payload.Type,
		Column1: payload.UserID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	// Extract notification IDs into a slice of int32
	ids := make([]int32, len(notifications))
	uids := make([]uuid.UUID, len(notifications)) // TODO: consider removing to just use payload.UserId
	for i, notification := range notifications {
		ids[i] = notification.ID
		uids[i] = notification.UserID
	}

	err = qtx.CreateNotificationPrayerInvite(ctx, sqlc.CreateNotificationPrayerInviteParams{
		AuthorID: payload.AuthorID,
		RoomID:   payload.RoomID,
		Column1:  uids,
		Column2:  ids,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// -- Notifications
// -- 1 user
// -- comment on their post(postComment), like their post(postLike), follow(follow), join church(churchJoin), repost(repost), deleted or blocked post.

// -- Many users
// -- invite to prayer(prayerInvite), church announcement(churchAnnouncement), mentioned in post(postMention),
