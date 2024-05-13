package services

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/steve-mir/diivix_backend/worker"
)

func SendEmail(taskDistributor worker.TaskDistributor, ctx context.Context, email, content string) error {
	log.Printf("Sending email to %s. Content %s", email, content)

	taskPayload := &worker.PayloadSendEmail{
		Username: email,
		Content:  content,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err := taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		return err //status.Errorf(codes.Internal, "failed to distribute task to send verify email %s", err)
	}

	return nil
}
