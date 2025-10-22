package redis

import (
	"fmt"
	"radius-server/src/config"
	"radius-server/src/database/entities"
	numberUtil "radius-server/src/utils/number"
	redisUtil "radius-server/src/utils/redis"
)

func GetNASByIP(ip string) (*entities.RadiusNas, error) {
	query := fmt.Sprintf("@ip_address:{%s}", redisUtil.PrepareParam(ip))

	res, err := redisClient.Do(Ctx, "FT.SEARCH", config.NasIndex, query, "LIMIT", "0", "1").Result()
	if err != nil {
		return nil, fmt.Errorf("redis search failed: %w", err)
	}

	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := resMap["total_results"].(int64)
	if !ok || totalResults == 0 {
		return nil, fmt.Errorf("record not found for IP: %s", ip)
	}

	results, ok := resMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("no results found for IP: %s", ip)
	}

	firstResult, ok := results[0].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid result format")
	}

	extraAttrs, ok := firstResult["extra_attributes"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid extra_attributes format")
	}

	fieldMap := make(map[string]string)
	for k, v := range extraAttrs {
		key := fmt.Sprintf("%v", k)
		val := fmt.Sprintf("%v", v)
		fieldMap[key] = val
	}

	nas := &entities.RadiusNas{}
	if id, ok := fieldMap["id"]; ok {
		value, _ := numberUtil.ParseToInt64(id)
		nas.Id = value
	}
	if name, ok := fieldMap["nas_name"]; ok {
		nas.NasName = &name
	}
	if ipAddr, ok := fieldMap["ip_address"]; ok {
		nas.IpAddress = ipAddr
	}
	if secret, ok := fieldMap["secret"]; ok {
		nas.Secret = secret
	}

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

	fields := map[string]interface{}{
		"id":         nas.Id,
		"ip_address": nas.IpAddress,
		"secret":     nas.Secret,
	}

	if nas.NasName != nil {
		fields["nas_name"] = *nas.NasName
	}
	if nas.SubscriberId != nil {
		fields["subscriber_id"] = *nas.SubscriberId
	}
	if nas.SessionId != nil {
		fields["session_id"] = *nas.SessionId
	}
	if nas.CreatedAt != 0 {
		fields["created_at"] = nas.CreatedAt
	}
	if nas.UpdatedAt != 0 {
		fields["updated_at"] = nas.UpdatedAt
	}

	if err := redisClient.HSet(Ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("failed to set NAS client in Redis: %w", err)
	}

	return nil
}

type SubscriberData struct {
	SubscriberID    string `json:"subscriber_id"`
	IP              string `json:"ip"`
	IpVersion       string `json:"ip_version"`
	SessionID       string `json:"session_id"`
	LastUpdatedTime int64  `json:"last_updated_time"`
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

func GetSubscriberByIP(ip string) (*SubscriberData, error) {
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

	subscriber := &SubscriberData{}
	if id, ok := fieldMap["subscriber_id"]; ok {
		subscriber.SubscriberID = id
	}
	if ipAddr, ok := fieldMap["ip"]; ok {
		subscriber.IP = ipAddr
	}
	if sessionID, ok := fieldMap["session_id"]; ok {
		subscriber.SessionID = sessionID
	}
	if lastUpdated, ok := fieldMap["last_updated_time"]; ok {
		value, _ := numberUtil.ParseToInt64(lastUpdated)
		subscriber.LastUpdatedTime = value
	}
	if ipVersion, ok := fieldMap["ip_version"]; ok {
		subscriber.IpVersion = ipVersion
	}

	return subscriber, nil
}

func GetSubscriberBySessionID(sessionID string) ([]*SubscriberData, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	query := fmt.Sprintf("@session_id:{%s}", redisUtil.PrepareParam(sessionID))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", config.SubscriberIndex, query).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search subscriber by session ID: %w", err)
	}

	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := resMap["total_results"].(int64)
	if !ok || totalResults == 0 {
		return nil, fmt.Errorf("subscriber not found for session ID: %s", sessionID)
	}

	results, ok := resMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("no results found for session ID: %s", sessionID)
	}

	var subscribers []*SubscriberData
	for _, result := range results {
		resultMap, ok := result.(map[interface{}]interface{})
		if !ok {
			continue
		}

		extraAttrs, ok := resultMap["extra_attributes"].(map[interface{}]interface{})
		if !ok {
			continue
		}

		fieldMap := make(map[string]string)
		for k, v := range extraAttrs {
			key := fmt.Sprintf("%v", k)
			val := fmt.Sprintf("%v", v)
			fieldMap[key] = val
		}

		subscriber := &SubscriberData{}
		if id, ok := fieldMap["subscriber_id"]; ok {
			subscriber.SubscriberID = id
		}
		if ip, ok := fieldMap["ip"]; ok {
			subscriber.IP = ip
		}
		if sessID, ok := fieldMap["session_id"]; ok {
			subscriber.SessionID = sessID
		}
		if lastUpdated, ok := fieldMap["last_updated_time"]; ok {
			value, _ := numberUtil.ParseToInt64(lastUpdated)
			subscriber.LastUpdatedTime = value
		}
		if ipVersion, ok := fieldMap["ip_version"]; ok {
			subscriber.IpVersion = ipVersion
		}

		subscribers = append(subscribers, subscriber)
	}

	return subscribers, nil
}

func CreateOrUpdateSubscriber(subscriber *SubscriberData) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber data cannot be nil")
	}

	if subscriber.IP == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	subscriberKey := fmt.Sprintf("%s:%s", config.SubscriberHashTableName, subscriber.IP)
	subscriberFields := map[string]interface{}{
		"subscriber_id":     subscriber.SubscriberID,
		"ip":                subscriber.IP,
		"ip_version":        subscriber.IpVersion,
		"session_id":        subscriber.SessionID,
		"last_updated_time": subscriber.LastUpdatedTime,
	}

	if err := redisClient.HSet(Ctx, subscriberKey, subscriberFields).Err(); err != nil {
		return fmt.Errorf("failed to set subscriber hash: %w", err)
	}

	return nil
}
