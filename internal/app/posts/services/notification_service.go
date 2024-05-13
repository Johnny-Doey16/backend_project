package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/steve-mir/diivix_backend/worker"
)

func SendNotification(taskDistributor worker.TaskDistributor, ctx context.Context, authorId, postId uuid.UUID, uids []uuid.UUID, notificationType, title, imageUrl, content, roomID string) error {
	taskPayload := &worker.PayloadSendNotification{
		UserID:   uids,
		Type:     notificationType,
		Title:    title,
		Body:     content,
		ImageUrl: imageUrl,
		PostID:   postId,
		AuthorID: authorId,
		RoomID:   roomID,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(5 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err := taskDistributor.DistributeTaskSendNotification(ctx, taskPayload, opts...)
	if err != nil {
		return err //status.Errorf(codes.Internal, "failed to distribute task to send verify email %s", err)
	}
	return nil
}
