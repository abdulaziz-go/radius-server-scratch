package redis

import (
	"fmt"
	"net"
	"radius-server/src/database/entities"
	numberUtil "radius-server/src/utils/number"
	redisUtil "radius-server/src/utils/redis"
)

func GetNASByIP(ip string) (*entities.RadiusNas, error) {
	query := fmt.Sprintf("@ip_address:{%s}", redisUtil.PrepareParam(ip))

	res, err := redisClient.Do(Ctx, "FT.SEARCH", nasIndex, query, "LIMIT", "0", "1").Result()
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

func HSetNasClient(fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields provided for NAS client")
	}

	id, ok := fields["id"]
	if !ok {
		return fmt.Errorf("missing 'id' field for NAS client")
	}

	key := fmt.Sprintf("%v:%v", nasHashTableName, id)

	args := make([]interface{}, 0, len(fields)*2+2)
	args = append(args, key)
	for k, v := range fields {
		args = append(args, k, v)
	}

	cmd := redisClient.Do(Ctx, append([]interface{}{"HSET"}, args...)...)
	if cmd.Err() != nil {
		return fmt.Errorf("failed to set NAS client in Redis: %w", cmd.Err())
	}

	return nil
}

type SubscriberData struct {
	SubscriberID    int64  `json:"subscriber_id"`
	IP              string `json:"ip"`
	SessionID       string `json:"session_id"`
	LastUpdatedTime int64  `json:"last_updated_time"`
}

func DeleteSubscriberByIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}

	query := fmt.Sprintf("@ip:{%s}", redisUtil.PrepareParam(ip))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", subscriberIndex, query, "LIMIT", "0", "1").Result()
	if err != nil {
		return fmt.Errorf("failed to search subscriber by IP: %w", err)
	}

	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := resMap["total_results"].(int64)
	if !ok || totalResults == 0 {
		return nil
	}

	results, ok := resMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil
	}

	firstResult, ok := results[0].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("invalid result format")
	}

	docKey, ok := firstResult["id"].(string)
	if !ok {
		return fmt.Errorf("invalid document key format")
	}

	if err := redisClient.Del(Ctx, docKey).Err(); err != nil {
		return fmt.Errorf("failed to delete subscriber by IP %s: %w", ip, err)
	}

	return nil
}

func DeleteSubscriberBySessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	query := fmt.Sprintf("@session_id:{%s}", redisUtil.PrepareParam(sessionID))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", subscriberIndex, query).Result()
	if err != nil {
		return fmt.Errorf("failed to search subscriber by session ID: %w", err)
	}

	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := resMap["total_results"].(int64)
	if !ok || totalResults == 0 {
		return nil
	}

	results, ok := resMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil
	}

	for _, result := range results {
		resultMap, ok := result.(map[interface{}]interface{})
		if !ok {
			continue
		}
		docKey, ok := resultMap["id"].(string)
		if !ok {
			continue
		}
		if err := redisClient.Del(Ctx, docKey).Err(); err != nil {
			return fmt.Errorf("failed to delete subscriber key %s: %w", docKey, err)
		}
	}

	return nil
}

func GetSubscriberByIP(ip string) (*SubscriberData, error) {
	if ip == "" {
		return nil, fmt.Errorf("IP address cannot be empty")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	query := fmt.Sprintf("@ip:{%s}", redisUtil.PrepareParam(ip))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", subscriberIndex, query, "LIMIT", "0", "1").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search subscriber by IP: %w", err)
	}
	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := resMap["total_results"].(int64)
	if !ok || totalResults == 0 {
		return nil, fmt.Errorf("subscriber not found for IP: %s", ip)
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

	subscriber := &SubscriberData{}
	if id, ok := fieldMap["subscriber_id"]; ok {
		value, _ := numberUtil.ParseToInt64(id)
		subscriber.SubscriberID = value
	}
	if ip, ok := fieldMap["ip"]; ok {
		subscriber.IP = ip
	}
	if sessionID, ok := fieldMap["session_id"]; ok {
		subscriber.SessionID = sessionID
	}
	if lastUpdated, ok := fieldMap["last_updated_time"]; ok {
		value, _ := numberUtil.ParseToInt64(lastUpdated)
		subscriber.LastUpdatedTime = value
	}

	return subscriber, nil
}

func GetSubscriberBySessionID(sessionID string) (*SubscriberData, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	query := fmt.Sprintf("@session_id:{%s}", redisUtil.PrepareParam(sessionID))
	res, err := redisClient.Do(Ctx, "FT.SEARCH", subscriberIndex, query, "LIMIT", "0", "1").Result()
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

	subscriber := &SubscriberData{}
	if id, ok := fieldMap["subscriber_id"]; ok {
		value, _ := numberUtil.ParseToInt64(id)
		subscriber.SubscriberID = value
	}
	if ip, ok := fieldMap["ip"]; ok {
		subscriber.IP = ip
	}
	if sessionID, ok := fieldMap["session_id"]; ok {
		subscriber.SessionID = sessionID
	}
	if lastUpdated, ok := fieldMap["last_updated_time"]; ok {
		value, _ := numberUtil.ParseToInt64(lastUpdated)
		subscriber.LastUpdatedTime = value
	}

	return subscriber, nil
}

func CreateOrUpdateSubscriber(subscriber *SubscriberData) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber data cannot be nil")
	}

	if subscriber.SubscriberID == 0 {
		return fmt.Errorf("subscriber ID cannot be empty")
	}

	subscriberKey := fmt.Sprintf("%s:%d", subscriberHashTableName, subscriber.SubscriberID)
	subscriberFields := map[string]interface{}{
		"subscriber_id":     subscriber.SubscriberID,
		"ip":                subscriber.IP,
		"session_id":        subscriber.SessionID,
		"last_updated_time": subscriber.LastUpdatedTime,
	}

	if err := redisClient.HSet(Ctx, subscriberKey, subscriberFields).Err(); err != nil {
		return fmt.Errorf("failed to set subscriber hash: %w", err)
	}

	return nil
}
