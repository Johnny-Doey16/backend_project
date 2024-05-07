package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	// "github.com/steve-mir/diivix_backend/internal/app/posts/pb"
)

type Cache struct {
	client *redis.Client
}

func NewCache(addr, username, password string, db int) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			Username: username,
		}),
	}
}

func (c *Cache) XAdd(ctx context.Context, stream, id string, values map[string]interface{}) (string, error) {
	args := &redis.XAddArgs{
		Stream:     stream,
		ID:         id, // You can use "*" to let Redis generate a unique ID
		Values:     values,
		MaxLen:     0,     // If you want to limit the length of the stream, set this to a positive number
		Approx:     true,  // Set to true to allow for approximate trimming of the stream to MaxLen
		NoMkStream: false, // Set to true if you don't want to create the stream if it doesn't exist
	}
	return c.client.XAdd(ctx, args).Result()
}

func (c *Cache) SetKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) XRead(ctx context.Context, streams []string, count int64, block time.Duration) ([]redis.XStream, error) {
	return c.client.XRead(ctx, &redis.XReadArgs{
		Streams: streams, // Use the format []string{"streamName", "lastID", ...}
		Count:   count,
		Block:   block,
	}).Result()
}

func (c *Cache) GetKey(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

func (c *Cache) XGroupCreate(ctx context.Context, stream, group, start string) error {
	return c.client.XGroupCreateMkStream(ctx, stream, group, start).Err()
}

func (c *Cache) XAck(ctx context.Context, stream, group, id string) (int64, error) {
	return c.client.XAck(ctx, stream, group, id).Result()
}

func (c *Cache) DeleteKey(ctx context.Context, key string) (int64, error) {
	return c.client.Del(ctx, key).Result()
}

// SIsMember checks if a member is in a Redis set.
func (c *Cache) SIsMember(ctx context.Context, setKey string, member interface{}) (bool, error) {
	cmd := c.client.SIsMember(ctx, setKey, member)
	return cmd.Result()
}

// SAdd adds a member to a Redis set.
func (c *Cache) SAdd(ctx context.Context, setKey string, member interface{}) error {
	cmd := c.client.SAdd(ctx, setKey, member)
	_, err := cmd.Result()
	return err
}

// SRem removes a member from a Redis set.
func (c *Cache) SRem(ctx context.Context, setKey string, member interface{}) error {
	cmd := c.client.SRem(ctx, setKey, member)
	_, err := cmd.Result()
	return err
}

// Worker function that listens to the channel and processes messages
// TODO: Refactor to be like PostSubscribeAndProcess and more generic
// func (c *Cache) SubscribeAndProcessOld(ctx context.Context, channel string) {
// 	pubsub := c.client.Subscribe(ctx, channel)
// 	defer pubsub.Close()

// 	// Wait for confirmation that subscription is created before starting to receive messages
// 	_, err := pubsub.Receive(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Go channel which receives messages
// 	ch := pubsub.Channel()

// 	for msg := range ch {
// 		fmt.Printf("Received message from %s: %s\n", msg.Channel, msg.Payload)
// 		// Process message
// 	}
// }

func (c *Cache) Publish(ctx context.Context, channel string, message []byte) *redis.IntCmd {
	return c.client.Publish(ctx, channel, message)
}

// func (c *Cache) SubscribeAndProcess(ctx context.Context, channel string) {
// 	pubsub := c.client.Subscribe(ctx, channel)
// 	defer pubsub.Close()

// 	// Wait for confirmation that subscription is created
// 	if _, err := pubsub.Receive(ctx); err != nil {
// 		log.Printf("Subscribe failed: %v", err)
// 		return // Exit the function without exiting the entire program
// 	}

// 	ch := pubsub.Channel()

// 	for {
// 		select {
// 		case msg := <-ch:
// 			// Process the received message
// 			var post pb.Post
// 			if err := json.Unmarshal([]byte(msg.Payload), &post); err != nil {
// 				log.Printf("Error unmarshalling message: %v", err)
// 				continue
// 			}
// 			// Handle the post object
// 		case <-ctx.Done():
// 			// Context was cancelled, stop processing
// 			log.Printf("Subscription to channel %s stopped", channel)
// 			return
// 		}
// 	}
// }

func (c *Cache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

// Define a struct for your message or task
type Task struct {
	Stream string
	Group  string
	ID     string
	Data   map[string]interface{}
}

// Worker function to process tasks
func worker(id int, tasks <-chan Task, cache *Cache) {
	for task := range tasks {
		// Process the task using task.Data
		// ...

		// After processing, acknowledge the message
		_, err := cache.XAck(context.Background(), task.Stream, task.Group, task.ID)
		if err != nil {
			// Handle error
		}
	}
}

// Function to create and start workers
func startWorkers(cache *Cache, numWorkers int) <-chan Task {
	tasks := make(chan Task)

	for i := 0; i < numWorkers; i++ {
		go worker(i, tasks, cache)
	}

	return tasks
}

/**
func main() {
    // Initialize a new cache with Redis client.
    cache := NewCache("localhost:6379", "", 0)

    // Set a key-value pair in the cache with an expiration time.
    err := cache.SetKey("foo", "bar", 0)
    if err != nil {
        panic(err)
    }

    // Retrieve the value for a given key.
    val, err := cache.GetKey("foo")
    if err != nil {
        panic(err)
    }
    fmt.Println("foo:", val)

    // Optionally, delete the key from the cache.
    _, err = cache.DeleteKey("foo")
    if err != nil {
        panic(err)
    }
}
*/

/**
var ctx = context.Background()

type Cache struct {
    client *redis.Client
}

// ... Cache struct methods including NewCache ...

func main() {
    // Initialize a new cache with Redis client.
    cache := NewCache("localhost:6379", "", 0)

    // Add message to a stream
    id, err := cache.XAdd("mystream", "*", map[string]interface{}{"key": "value"})
    if err != nil {
        panic(err)
    }
    fmt.Println("Added message ID:", id)

    // Read messages from a stream
    messages, err := cache.XRead([]string{"mystream", "0"}, 2, 0)
    if err != nil {
        panic(err)
    }
    for _, stream := range messages {
        for _, message := range stream.Messages {
            fmt.Printf("Message ID: %s, Values: %v\n", message.ID, message.Values)
        }
    }

    // Create a consumer group (if it doesn't exist)
    err = cache.XGroupCreate("mystream", "mygroup", "$")
    if err != nil {
        panic(err)
    }

    // Read messages using XREADGROUP (implementing consumer logic is required)
    // Acknowledge messages using XACK after processing
    // ...
}
*/
