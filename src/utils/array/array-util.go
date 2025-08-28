package arrayUtil

import "fmt"

func Filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, entity := range items {
		if predicate(entity) {
			result = append(result, entity)
		}
	}
	return result
}

func FindItem[T any](items []T, conditions []func(T) bool) *T {
	var result *T
	for _, entity := range items {
		meetsAllConditions := true
		for _, condition := range conditions {
			if !condition(entity) {
				meetsAllConditions = false
				break
			}
		}
		if meetsAllConditions {
			result = &entity
			break
		}
	}
	return result
}

func ItemExists[T comparable](array []T, item T) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

func UniqueInt64(slice []int64) []int64 {
	uniqueMap := make(map[int64]bool)
	var uniqueSlice []int64
	for _, v := range slice {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			uniqueSlice = append(uniqueSlice, v)
		}
	}
	return uniqueSlice
}

func DerefSlice[T any](array []*T) []T {
	if len(array) == 0 {
		return make([]T, 0)
	}
	slice := make([]T, 0)
	for _, item := range array {
		if item != nil {
			slice = append(slice, *item)
		}
	}
	return slice
}

func SliceArrayFromValue(arr []string, value string) []string {
	for i, v := range arr {
		if v == value {
			return arr[i:]
		}
	}
	return nil
}

func SliceArrayBeforeValue(arr []string, value string) []string {
	for i, v := range arr {
		if v == value {
			return arr[:i+1]
		}
	}
	return arr
}

func RemoveStrings(slice []string, toRemove []string) []string {
	result := []string{}
	removeMap := make(map[string]bool)
	for _, val := range toRemove {
		removeMap[val] = true
	}
	for _, item := range slice {
		if !removeMap[item] {
			result = append(result, item)
		}
	}
	return result
}

func EnumsToStrings[T ~string](items []T) []string {
	result := make([]string, 0)
	for _, item := range items {
		result = append(result, fmt.Sprintf("%s", item))
	}
	return result
}

func SliceArray[T any](arr []T, offset int64, count int) []T {
	if offset >= int64(len(arr)) || count <= 0 {
		return []T{}
	}
	if offset+int64(count) > int64(len(arr)) {
		count = len(arr) - int(offset)
	}
	return arr[offset : offset+int64(count)]
}
