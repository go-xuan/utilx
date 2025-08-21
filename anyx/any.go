package anyx

// If if取值
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// IfZero 为空时取默认值
func IfZero[T any](zero, def T) T {
	if ValueOf(&zero).IsZero() {
		return def
	}
	return zero
}

// Default 用于函数中的不定参数取默认值,可变
func Default[T any](def T, variable ...T) T {
	if len(variable) == 0 {
		return def
	}
	return variable[0]
}
