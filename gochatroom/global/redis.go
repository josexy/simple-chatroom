package global

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Hostname   string   `yaml:"hostname"`
	Port       int      `yaml:"port"`
	Password   string   `yaml:"password"`
	DB         int      `yaml:"db"`
	Cluster    bool     `yaml:"cluster_enable"`
	CLusterIPs []string `yaml:"cluster_ips"`
	Addr       string
}

func InitRedis(config *RedisConfig) {
	if config.Cluster {
		Redis = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        config.CLusterIPs,
			Password:     config.Password,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		})
	} else {
		Redis = redis.NewClient(&redis.Options{
			Addr:         config.Addr,
			Password:     config.Password,
			DB:           config.DB,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		})
	}
	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
