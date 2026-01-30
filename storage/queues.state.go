package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func StoreJobData(client *redis.Client, queueName string, jobID string, payload []byte) {
	if client == nil {
		return
	}

	redisKey := fmt.Sprintf("queue:%s:job:%s", queueName, jobID)

	ctx := context.Background()
	client.Set(ctx, redisKey, payload, 0) // Adding the job data per job id
}

func PushToActive(client *redis.Client, queueName string, jobID string) {
	if client == nil {
		return
	}

	ctx := context.Background()
	newActiveKey := fmt.Sprintf("queue:%s:active", queueName)

	client.LPush(ctx, newActiveKey, jobID)
}

// func PushToFinished(client *redis.Client, queueName string) (string, error) {
// 	if client == nil {
// 		return "", nil
// 	}

// 	oldKey := fmt.Sprintf("queue:%s:%s", queueName, JOB_STATUSES[1])
// 	newActiveKey := fmt.Sprintf("queue:%s:%s", queueName, JOB_STATUSES[2])

// 	ctx := context.Background()
// 	jobID, err := client.BRPopLPush(ctx, oldKey, newActiveKey, 0).Result()
// 	return jobID, err
// }

// func PushToFailed(client *redis.Client, queueName string) (string, error) {
// 	if client == nil {
// 		return "", nil
// 	}

// 	oldKey := fmt.Sprintf("queue:%s:%s", queueName, JOB_STATUSES[2])
// 	newActiveKey := fmt.Sprintf("queue:%s:%s", queueName, JOB_STATUSES[3])

// 	ctx := context.Background()
// 	jobID, err := client.BRPopLPush(ctx, oldKey, newActiveKey, 0).Result()
// 	return jobID, err
// }
