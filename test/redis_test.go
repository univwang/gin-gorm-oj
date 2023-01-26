package test

import (
	"context"
	"fmt"
	"gin_oj/models"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "39.101.1.158:3004",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestRedisSet(t *testing.T) {
	//set

	rdb.Set(ctx, "name", "mmc", time.Second*30)
}

func TestRedisGet(t *testing.T) {
	//get
	result, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestRedisGetByModels(t *testing.T) {
	result, err := models.RDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
