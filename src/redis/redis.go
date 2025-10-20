package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"radius-server/src/common/logger"
	"radius-server/src/config"
	timeUtil "radius-server/src/utils/time"
)

var (
	redisClient *redis.Client
	Ctx         = context.Background()
)

func Connect() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:            config.AppConfig.Redis.Host,
		Protocol:        2,
		PoolSize:        config.AppConfig.Redis.Connection.MaxNumber,
		MinIdleConns:    config.AppConfig.Redis.Connection.OpenMinNumber,
		ConnMaxLifetime: timeUtil.DurationSeconds(config.AppConfig.Redis.Connection.MaxLifetimeSec),
		ConnMaxIdleTime: timeUtil.DurationSeconds(config.AppConfig.Redis.Connection.MaxIdleTimeSec),
	})

	if err := redisClient.Ping(Ctx).Err(); err != nil {
		logger.Logger.Fatal().Msgf("Connection to redis error. %s", err.Error())
		return err
	}

	logger.Logger.Info().Msg("Redis Client connected successfully")

	return nil
}
