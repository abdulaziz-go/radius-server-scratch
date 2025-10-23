package redis

import (
	"fmt"
	"radius-server/src/config"
	"radius-server/src/database/entities"
	redisUtil "radius-server/src/utils/redis"
)

func GetNASByIP(ip string) (*entities.RadiusNas, error) {
	query := fmt.Sprintf("@ip_address:{%s}", redisUtil.PrepareParam(ip))

	res, err := redisClient.Do(Ctx, "FT.SEARCH", config.NasIndex, query, "LIMIT", "0", "1").Result()
	if err != nil {
		return nil, fmt.Errorf("redis search failed: %w", err)
	}

	fieldMap, err := redisUtil.ParseSingleRedisSearchResult(res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search result for IP %s: %w", ip, err)
	}

	nas := &entities.RadiusNas{}
	nas = nas.NewRadiusNasFromFieldMap(fieldMap)
	return nas, nil
}

func HSetNasClient(nas *entities.RadiusNas) error {
	if nas == nil {
		return fmt.Errorf("NAS client cannot be nil")
	}

	if nas.Id == 0 {
		return fmt.Errorf("missing 'id' field for NAS client")
	}

	key := fmt.Sprintf("%v:%v", config.NasHashTableName, nas.Id)
	fields := nas.ToFieldsMap()

	if err := redisClient.HSet(Ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("failed to set NAS client in Redis: %w", err)
	}

	return nil
}

func DeleteSubscriberByIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	key := fmt.Sprintf("%s:%s", config.SubscriberHashTableName, ip)

	if err := redisClient.Del(Ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete subscriber by IP %s: %w", ip, err)
	}

	return nil
}

func GetSubscriberByIP(ip string) (*entities.SubscriberData, error) {
	if ip == "" {
		return nil, fmt.Errorf("IP address cannot be empty")
	}

	key := fmt.Sprintf("%s:%s", config.SubscriberHashTableName, ip)

	fieldMap, err := redisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriber by IP: %w", err)
	}

	if len(fieldMap) == 0 {
		return nil, fmt.Errorf("subscriber not found for IP: %s", ip)
	}

	subscriber := &entities.SubscriberData{}
	subscriber = subscriber.NewSubscriberDataFromFieldMap(fieldMap)

	return subscriber, nil
}

func GetSubscriberBySessionID(sessionID string) ([]*entities.SubscriberData, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	query := fmt.Sprintf("@session_id:{%s}", redisUtil.PrepareParam(sessionID))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", config.SubscriberIndex, query).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search subscriber by session ID: %w", err)
	}

	results, err := redisUtil.ParseRedisSearchResponse(res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results for session ID %s: %w", sessionID, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("subscriber not found for session ID: %s", sessionID)
	}

	var subscribers []*entities.SubscriberData
	for _, fieldMap := range results {
		subscriber := &entities.SubscriberData{}
		subscriber = subscriber.NewSubscriberDataFromFieldMap(fieldMap)
		subscribers = append(subscribers, subscriber)
	}

	return subscribers, nil
}

func CreateOrUpdateSubscriber(subscriber *entities.SubscriberData) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber data cannot be nil")
	}

	if subscriber.IP == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	subscriberKey := fmt.Sprintf("%s:%s", config.SubscriberHashTableName, subscriber.IP)
	subscriberFields := subscriber.ToFieldsMap()

	if err := redisClient.HSet(Ctx, subscriberKey, subscriberFields).Err(); err != nil {
		return fmt.Errorf("failed to set subscriber hash: %w", err)
	}

	return nil
}

func GetSubscriberBySubscriberID(subscriberID string) ([]*entities.SubscriberData, error) {
	if subscriberID == "" {
		return nil, fmt.Errorf("subscriber ID cannot be empty")
	}

	query := fmt.Sprintf("@subscriber_id:{%s}", redisUtil.PrepareParam(subscriberID))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", config.SubscriberIndex, query).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search subscriber by subscriber ID: %w", err)
	}

	results, err := redisUtil.ParseRedisSearchResponse(res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results for subscriber ID %s: %w", subscriberID, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("subscriber not found for subscriber ID: %s", subscriberID)
	}

	var subscribers []*entities.SubscriberData
	for _, fieldMap := range results {
		subscriber := &entities.SubscriberData{}
		subscriber = subscriber.NewSubscriberDataFromFieldMap(fieldMap)
		subscribers = append(subscribers, subscriber)
	}

	return subscribers, nil
}
