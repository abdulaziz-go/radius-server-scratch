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

<<<<<<< HEAD
	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := resMap["total_results"].(int64)
=======
	arr, ok := res.([]interface{})
	if !ok || len(arr) < 2 {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	totalResults, ok := arr[0].(int64)
>>>>>>> dbc31cde1d97d7490e03beb6da360df823da028e
	if !ok || totalResults == 0 {
		return nil, fmt.Errorf("record not found for IP: %s", ip)
	}

	if len(arr) < 3 {
		return nil, fmt.Errorf("no results found for IP: %s", ip)
	}

	// Fields array for the first document
	docFields, ok := arr[2].([]interface{}) // arr[1] is document ID
	if !ok {
		return nil, fmt.Errorf("invalid document format")
	}

<<<<<<< HEAD
	extraAttrs, ok := firstResult["extra_attributes"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid extra_attributes format")
	}

=======
	// Convert flat array to map[string]string
>>>>>>> dbc31cde1d97d7490e03beb6da360df823da028e
	fieldMap := make(map[string]string)
	for i := 0; i < len(docFields)-1; i += 2 {
		key, _ := docFields[i].(string)
		val := fmt.Sprintf("%v", docFields[i+1])
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

	// Execute HSET command
	cmd := redisClient.Do(Ctx, append([]interface{}{"HSET"}, args...)...)
	if cmd.Err() != nil {
		return fmt.Errorf("failed to set NAS client in Redis: %w", cmd.Err())
	}

	return nil
}
