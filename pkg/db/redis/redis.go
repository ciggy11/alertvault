package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/ciggy11/alertvault/pkg/alert"
)

type RedisConfig struct {
	Addrs        string        `yaml:"addrs" env:"REDIS_ADDRS" default:"0.0.0.0:6379"`
	Password     string        `yaml:"password" env:"REDIS_PASSWORD"`
	AlertsDB     int           `yaml:"alerts_db" env:"ALERTS_DB" default:"0"`
	AlertGroupDB int           `yaml:"alert_group_db" env:"AlertGroup_DB" default:"1"`
	MasterName   string        `yaml:"master_name" env:"REDIS_MASTER_NAME"`
	Timeout      time.Duration `yaml:"timeout"`
	Expiration   time.Duration `yaml:"expiration"`
}

type RedisClient struct {
	expiration time.Duration
	timeout    time.Duration
	rdb        goredis.UniversalClient
}



func NewRedisAlertsClient(cfg *RedisConfig) *RedisClient {
	opt := goredis.UniversalOptions{
		Addrs:      strings.Split(cfg.Addrs, ","),
		MasterName: cfg.MasterName,
		Password:   cfg.Password,
		DB:         cfg.AlertsDB,
	}
	client := &RedisClient{
		expiration: cfg.Expiration,
		timeout:    cfg.Timeout,
		rdb:        goredis.NewUniversalClient(&opt),
	}
	err := client.Ping(context.Background())
	if err != nil {
		log.Errorf("Failed to connect to redis: %s", err)
		panic(err)
	}
	return client
}

func NewRedisGroupAlertsClient(cfg *RedisConfig) *RedisClient {
	opt := goredis.UniversalOptions{
		Addrs:      strings.Split(cfg.Addrs, ","),
		MasterName: cfg.MasterName,
		Password:   cfg.Password,
		DB:         cfg.AlertGroupDB,
	}
	client := &RedisClient{
		expiration: cfg.Expiration,
		timeout:    cfg.Timeout,
		rdb:        goredis.NewUniversalClient(&opt),
	}
	err := client.Ping(context.Background())
	if err != nil {
		log.Errorf("Failed to connect to redis: %s", err)
		panic(err)
	}
	return client
}

func (c *RedisClient) Ping(ctx context.Context) error {
	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	pong, err := c.rdb.Ping(ctx).Result()

	if err != nil {
		return err
	}

	if pong != "PONG" {
		return fmt.Errorf("redis: Unexpected Ping response %q", pong)
	}
	return nil
}

func (c *RedisClient) ZSet(ctx context.Context, member []byte, key string, score float64) error {
	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	_, err := c.rdb.ZAdd(ctx, key, &goredis.Z{
		Score:  score,
		Member: member,
	}).Result()
	if err != nil {
		log.Errorf("Failed to set alert: %s", err)
		panic(err)
	}
	return err
}
func (c *RedisClient) ZGetByScore(ctx context.Context, ad *alert.AlertsDesc) ([]*alert.Alert, error) {
	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	result, err := c.rdb.ZRangeByScore(ctx, ad.Key, &goredis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatFloat(ad.Score, 'f', -1, 64),
		Offset: ad.Offset,
		Count:  ad.Count,
	}).Result()
	if err != nil {
		return nil, err
	}
	tenantAlertData, err := alert.StringToAlerts(result)
	if err != nil {
		return nil, err
	}
	return tenantAlertData, nil
}
func (c *RedisClient) ZGetAll(ctx context.Context, key string) ([]string, error) {
	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	result, err := c.rdb.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RedisClient) Set(ctx context.Context, key string, value []byte) error {
	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	err := c.rdb.Set(ctx, key, value, c.expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c *RedisClient) Exist(ctx context.Context, key string) (bool, error) {
	val, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

func (c *RedisClient) Count(ctx context.Context, key string) (int64, error) {
	val, err := c.rdb.ZCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}
