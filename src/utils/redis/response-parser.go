package redisUtil

import (
	"fmt"
	"strconv"
)

func ParseRedisSearchResponse(res interface{}) ([]map[string]string, error) {
	switch v := res.(type) {
	case []interface{}:
		return parseArrayResponse(v)
	case map[interface{}]interface{}:
		return parseMapResponse(v)
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

func ParseSingleRedisSearchResult(res interface{}) (map[string]string, error) {
	results, err := ParseRedisSearchResponse(res)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	return results[0], nil
}

func parseArrayResponse(arr []interface{}) ([]map[string]string, error) {
	if len(arr) < 1 {
		return nil, fmt.Errorf("empty response array")
	}

	countStr := fmt.Sprintf("%v", arr[0])
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, fmt.Errorf("invalid count format: %v", arr[0])
	}

	if count == 0 {
		return []map[string]string{}, nil
	}

	var results []map[string]string

	for i := 1; i < len(arr); i += 2 {
		if i+1 >= len(arr) {
			break
		}

		fieldsInterface := arr[i+1]
		fieldsArray, ok := fieldsInterface.([]interface{})
		if !ok {
			continue
		}

		fieldMap := make(map[string]string)
		for j := 0; j < len(fieldsArray); j += 2 {
			if j+1 < len(fieldsArray) {
				key := fmt.Sprintf("%v", fieldsArray[j])
				value := fmt.Sprintf("%v", fieldsArray[j+1])
				fieldMap[key] = value
			}
		}

		if len(fieldMap) > 0 {
			results = append(results, fieldMap)
		}
	}

	return results, nil
}

func parseMapResponse(resMap map[interface{}]interface{}) ([]map[string]string, error) {
	totalResults, ok := resMap["total_results"].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid total_results format")
	}

	if totalResults == 0 {
		return []map[string]string{}, nil
	}

	results, ok := resMap["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid results format")
	}

	var parsedResults []map[string]string
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

		parsedResults = append(parsedResults, fieldMap)
	}

	return parsedResults, nil
}
