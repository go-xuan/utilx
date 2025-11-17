package slicex

import "github.com/go-xuan/typex"

// Conv2Map 数组转map
func Conv2Map[T typex.Unique](slice []T) map[string]T {
	if len(slice) == 0 {
		return nil
	}
	var m = make(map[string]T)
	for _, item := range slice {
		key := item.GetKey()
		m[key] = item
	}
	return m
}

// Conv2GroupMap 数组分组
func Conv2GroupMap[T typex.Unique](slice []T) map[string][]T {
	if len(slice) == 0 {
		return nil
	}
	var groups = make(map[string][]T)
	for _, item := range slice {
		key := item.GetKey()
		groups[key] = append(groups[key], item)
	}
	return groups
}
