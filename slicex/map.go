package slicex

type Value interface {
	GetPrimaryKey() string // 获取唯一主键
}

// Conv2Map 数组转map
func Conv2Map[V Value](slice []V) map[string]V {
	if len(slice) == 0 {
		return nil
	}
	var m = make(map[string]V)
	for _, item := range slice {
		key := item.GetPrimaryKey()
		m[key] = item
	}
	return m
}

// Conv2Groups 数组分组
func Conv2Groups[V Value](slice []V) map[string][]V {
	if len(slice) == 0 {
		return nil
	}
	var groups = make(map[string][]V)
	for _, item := range slice {
		key := item.GetPrimaryKey()
		groups[key] = append(groups[key], item)
	}
	return groups
}
