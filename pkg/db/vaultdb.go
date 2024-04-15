package db

import (
	"github.com/ciggy11/alertvault/pkg/db/redis"
	"github.com/ciggy11/alertvault/pkg/db/s3"
)
type Config struct {
	Redis redis.RedisConfig `yaml:"redis"`
	S3    s3.S3Config       `yaml:"s3"`
}