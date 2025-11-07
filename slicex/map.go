package slicex

type IMap interface {
	MapKey() string
}

// Conv2Map 数组转map
func Conv2Map[T IMap](slice []T) map[string]T {
	if len(slice) == 0 {
		return nil
	}
	var m = make(map[string]T)
	for _, item := range slice {
		key := item.MapKey()
		m[key] = item
	}
	return m
}

// Conv2GroupMap 数组分组
func Conv2GroupMap[T IMap](slice []T) map[string][]T {
	if len(slice) == 0 {
		return nil
	}
	var groups = make(map[string][]T)
	for _, item := range slice {
		key := item.MapKey()
		groups[key] = append(groups[key], item)
	}
	return groups
}
