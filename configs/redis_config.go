package configs

import (
	"context"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisDB *redis.Client

func NewRedisConn() {
	config := commons.Configs.Datasource.Redis
	RedisDB = redis.NewClient(&redis.Options{
		Network:      config.Network,
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.DB,
		MinIdleConns: config.IdleMin,
		PoolSize:     config.PoolSize,
		PoolTimeout:  config.PoolTimeout * time.Minute,
	})
	err := RedisDB.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}
