package objectUtil

import (
	"encoding/json"
	"errors"
	"reflect"
)

func MarshalUnmarshal[T any](input any) (T, error) {
	var output T
	data, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	err = json.Unmarshal(data, &output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func GetJSONFields[T any]() []string {
	var fields []string
	var s T
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fields
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			fields = append(fields, jsonTag)
		}
	}
	return fields
}

func GetFieldNameByJSONTag(obj interface{}, jsonTag string) (*string, error) {
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonName := field.Tag.Get("json")
		if jsonName == jsonTag {
			return &field.Name, nil
		}
	}
	return nil, errors.New("Field not found")
}
