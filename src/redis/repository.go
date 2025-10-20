package redis

import (
	"fmt"
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

	// Parse map response
	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, _ := resMap["total_results"].(int64)
	if totalResults == 0 {
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
		return nil, fmt.Errorf("no extra_attributes in result")
	}

	// Convert to string map
	fieldMap := make(map[string]string)
	for k, v := range extraAttrs {
		key := fmt.Sprintf("%v", k)
		val := fmt.Sprintf("%v", v)
		fieldMap[key] = val
	}

	// Build NAS entity
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

	// Execute HSET command
	cmd := redisClient.Do(Ctx, append([]interface{}{"HSET"}, args...)...)
	if cmd.Err() != nil {
		return fmt.Errorf("failed to set NAS client in Redis: %w", cmd.Err())
	}

	return nil
}
