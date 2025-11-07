package slicex

// Contains 是否包含
func Contains[T comparable](slice []T, v T) bool {
	for _, item := range slice {
		if item == v {
			return true
		}
	}
	return false
}

// ContainsAny 是否包含任意元素
func ContainsAny[T comparable](slice []T, args ...T) bool {
	if args != nil && len(args) > 0 {
		m := make(map[T]struct{})
		for _, item := range slice {
			m[item] = struct{}{}
		}
		for _, arg := range args {
			if _, ok := m[arg]; ok {
				return true
			}
		}
	}
	return false
}

// ContainsAll 是否包含全部元素
func ContainsAll[T comparable](slice []T, args ...T) bool {
	if args != nil && len(args) > 0 {
		var m = make(map[T]struct{})
		for _, k := range slice {
			m[k] = struct{}{}
		}
		for _, arg := range args {
			if _, ok := m[arg]; !ok {
				return false
			}
		}
	}
	return true
}

// Distinct 去重
func Distinct[T comparable](slices ...[]T) []T {
	if len(slices) > 0 {
		var m = make(map[T]struct{})
		for _, slice := range slices {
			for _, k := range slice {
				m[k] = struct{}{}
			}
		}
		var set []T
		for k := range m {
			set = append(set, k)
		}
		return set
	}
	return nil
}

// Intersect 取交集
func Intersect[T comparable](slices ...[]T) []T {
	if l := len(slices); l > 0 {
		var m = make(map[T]int)
		for _, slice := range slices {
			for _, k := range slice {
				m[k]++
			}
		}
		var result []T
		for k, count := range m {
			if count == l {
				result = append(result, k)
			}
		}
		return result
	}
	return nil
}

// Exclude 移除
func Exclude[T comparable](target []T, exclude []T) []T {
	if len(target) > 0 && len(exclude) > 0 {
		var m = make(map[T]struct{})
		for _, k := range exclude {
			m[k] = struct{}{}
		}
		var result []T
		for _, k := range target {
			if _, ok := m[k]; !ok {
				result = append(result, k)
			}
		}
		return result
	}
	return target
}
