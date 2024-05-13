package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendEmail,
		opts ...asynq.Option,
	) error
	DistributeTaskSendNotification(
		ctx context.Context,
		payload *PayloadSendNotification,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {

	return &RedisTaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}
