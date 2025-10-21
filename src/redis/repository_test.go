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
		"SCHEMA", "ip_address", "TAG",
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

func TestSubscriberOperations(t *testing.T) {
	client := initTestRedis(t)
	client.FlushDB(testCtx)

	_, err := client.Do(testCtx,
		"FT.CREATE", subscriberIndex, "ON", "HASH", "PREFIX", "1", "subscriber:",
		"SCHEMA", "subscriber_id", "NUMERIC", "SORTABLE", "ip", "TEXT", "session_id", "TEXT", "last_updated_time", "NUMERIC", "SORTABLE",
	).Result()
	if err != nil && !strings.Contains(err.Error(), "Index already exists") {
		require.NoError(t, err, "Failed to create subscriber RediSearch index")
	}

	subscriber := &SubscriberData{
		SubscriberID:    12345,
		IP:              "192.168.1.100",
		SessionID:       "session-abc-123",
		LastUpdatedTime: 1634567890,
	}

	t.Run("Create subscriber", func(t *testing.T) {
		err := CreateOrUpdateSubscriber(subscriber)
		require.NoError(t, err, "CreateOrUpdateSubscriber should not return an error")

		stored, err := client.HGetAll(testCtx, fmt.Sprintf("%s:%d", subscriberHashTableName, subscriber.SubscriberID)).Result()
		require.NoError(t, err)
		assert.Equal(t, "12345", stored["subscriber_id"])
		assert.Equal(t, "192.168.1.100", stored["ip"])
		assert.Equal(t, "session-abc-123", stored["session_id"])
		assert.Equal(t, "1634567890", stored["last_updated_time"])
	})

	t.Run("Get subscriber by IP", func(t *testing.T) {
		retrieved, err := GetSubscriberByIP("192.168.1.100")
		require.NoError(t, err, "GetSubscriberByIP should not return an error")

		assert.Equal(t, int64(12345), retrieved.SubscriberID)
		assert.Equal(t, "192.168.1.100", retrieved.IP)
		assert.Equal(t, "session-abc-123", retrieved.SessionID)
		assert.Equal(t, int64(1634567890), retrieved.LastUpdatedTime)
	})

	t.Run("Get subscriber by session ID", func(t *testing.T) {
		retrieved, err := GetSubscriberBySessionID("session-abc-123")
		require.NoError(t, err, "GetSubscriberBySessionID should not return an error")

		assert.Equal(t, int64(12345), retrieved.SubscriberID)
		assert.Equal(t, "192.168.1.100", retrieved.IP)
		assert.Equal(t, "session-abc-123", retrieved.SessionID)
		assert.Equal(t, int64(1634567890), retrieved.LastUpdatedTime)
	})

	t.Run("Update subscriber", func(t *testing.T) {
		updatedSubscriber := &SubscriberData{
			SubscriberID:    12345,
			IP:              "192.168.1.200", // Changed IP
			SessionID:       "session-abc-123",
			LastUpdatedTime: 1634567900, // Updated time
		}

		err := CreateOrUpdateSubscriber(updatedSubscriber)
		require.NoError(t, err, "CreateOrUpdateSubscriber should not return an error")

		retrieved, err := GetSubscriberBySessionID("session-abc-123")
		require.NoError(t, err)
		assert.Equal(t, "192.168.1.200", retrieved.IP)
		assert.Equal(t, int64(1634567900), retrieved.LastUpdatedTime)
	})

	t.Run("Delete subscriber by IP", func(t *testing.T) {
		err := DeleteSubscriberByIP("192.168.1.200")
		require.NoError(t, err, "DeleteSubscriberByIP should not return an error")

		_, err = GetSubscriberByIP("192.168.1.200")
		assert.Error(t, err, "Should return error when subscriber not found")
	})

	t.Run("Delete subscriber by session ID", func(t *testing.T) {
		newSubscriber := &SubscriberData{
			SubscriberID:    54321,
			IP:              "192.168.1.150",
			SessionID:       "session-xyz-456",
			LastUpdatedTime: 1634567800,
		}
		err := CreateOrUpdateSubscriber(newSubscriber)
		require.NoError(t, err)

		err = DeleteSubscriberBySessionID("session-xyz-456")
		require.NoError(t, err, "DeleteSubscriberBySessionID should not return an error")

		_, err = GetSubscriberBySessionID("session-xyz-456")
		assert.Error(t, err, "Should return error when subscriber not found")
	})
}
