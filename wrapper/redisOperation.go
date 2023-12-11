package wrapper

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type LogObject struct {
	Sql       string    `json:"sql"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// MarshalBinary returns a JSON-encoded byte slice of the LogObject.
//
// Returns:
//   - data: the JSON-encoded byte slice of the LogObject
//   - err: any error that occurred during the marshaling process
//
// NOTE: without implementing it is giving error
func (l LogObject) MarshalBinary() (data []byte, err error) {
	return json.Marshal(l)
}

type RedisOperator interface {
	Push(ctx context.Context, key string, values ...interface{}) error
	Range(ctx context.Context, key string, start, stop int64) ([]string, error)
}

type redisOperator struct {
	client *redis.Client
}

// Push adds one or more values to the beginning of a list in Redis.
//
// Parameters:
//   - ctx: the context.Context used for the operation
//   - key: the key of the list in Redis
//   - values: one or more values to be pushed to the list, as interface{} types
//
// Returns:
//   - err: any error that occurred during the LPush operation
func (r *redisOperator) Push(ctx context.Context, key string, values ...interface{}) error {
	err := r.client.LPush(ctx, key, values).Err()
	return err
}

// Range returns a range of elements from a list in Redis.
//
// Parameters:
//   - ctx: the context.Context object for cancellation support
//   - key: the key of the list
//   - start: the starting index of the range
//   - stop: the ending index of the range
//
// start 0, stop -1 for getting all log
// Returns:
//   - []string: a slice of strings containing the elements within the specified range
//   - error: any error that occurred during the LRange operation
func (r *redisOperator) Range(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

func NewRedisOperator(client *redis.Client) RedisOperator {
	return &redisOperator{client: client}
}
