package bridge

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

var ctx = context.Background()
var rds *redis.Client
var rdsJson = rejson.NewReJSONHandler()

func getPostByToken(token string) (post Post, err error) {
	postJson, err := rdsJson.JSONGet(token, ".")
	if err != nil {
		return Post{}, err
	}
	err = json.Unmarshal(*postJson.(*[]byte), &post)
	return
}

func getPostByMediaGroupId(mediaGroupId string) (Post, error) {
	token, err := rds.Get(ctx, mediaGroupId).Result()
	if err != nil {
		return Post{}, err
	}
	return getPostByToken(token)
}

func initRedis() {
	rds = redis.NewClient(&redis.Options{})
	defer func() {
		if err := rds.FlushAll(ctx).Err(); err != nil {
			log.Fatalf("goredis - failed to flush: %v", err)
		}
		if err := rds.Close(); err != nil {
			log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
		}
	}()
	rdsJson.SetGoRedisClient(rds)
}
