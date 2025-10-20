package redis

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testRedisAddr = "localhost:6481"
	testRedisDB   = 0
	testCtx       = context.Background()
)

func initTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: testRedisAddr,
		DB:   testRedisDB,
	})
	err := client.Ping(testCtx).Err()
	require.NoError(t, err, "Redis must be running locally for tests")
	redisClient = client
	Ctx = testCtx
	return client
}

func TestHSetNasClientAndGetNASByIP(t *testing.T) {
	client := initTestRedis(t)
	client.FlushDB(testCtx)

	_, err := client.Do(testCtx,
		"FT.CREATE", nasIndex, "ON", "HASH", "PREFIX", "1", "radius_nas:",
		"SCHEMA", "id", "NUMERIC", "nas_name", "TEXT", "ip_address", "TAG", "secret", "TEXT",
	).Result()
	if err != nil && !strings.Contains(err.Error(), "Index already exists") {
		require.NoError(t, err, "Failed to create RediSearch index")
	}

	fields := map[string]interface{}{
		"id":         1,
		"nas_name":   "MainNAS",
		"ip_address": "192.168.1.10",
		"secret":     "testSecret",
	}

	t.Run("Insert NAS client with HSET", func(t *testing.T) {
		err := HSetNasClient(fields)
		require.NoError(t, err, "HSetNasClient should not return an error")

		stored, err := client.HGetAll(testCtx, fmt.Sprintf("%v:%v", nasHashTableName, fields["id"])).Result()
		require.NoError(t, err)
		assert.Equal(t, "MainNAS", stored["nas_name"])
		assert.Equal(t, "192.168.1.10", stored["ip_address"])
		assert.Equal(t, "testSecret", stored["secret"])
	})

	t.Run("Get NAS client by IP", func(t *testing.T) {
		nas, err := GetNASByIP("192.168.1.10")
		require.NoError(t, err, "GetNASByIP should not return an error")

		assert.Equal(t, int64(1), nas.Id)
		assert.Equal(t, "MainNAS", *nas.NasName)
		assert.Equal(t, "192.168.1.10", nas.IpAddress)
		assert.Equal(t, "testSecret", nas.Secret)
	})
}
