package mathx

import "math"

// Ratio 计算百分率
func Ratio(numerator, denominator float64, prec int) float64 {
	if denominator > 0 {
		return Ground(numerator*100/denominator, prec)
	}
	return 0
}

// Ground 四舍五入
func Ground(f float64, prec int) float64 {
	if f != 0 {
		pow := math.Pow10(prec)
		return math.Floor(f*pow+0.5) / pow
	}
	return 0
}

// Min 三数取小
func Min[T ~int | ~int32 | ~int64 | ~uint | float32 | ~float64](a, b, c T) T {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

// Max 三数取大
func Max[T ~int | ~int32 | ~int64 | ~uint | float32 | ~float64](a, b, c T) T {
	if a >= b && a >= c {
		return a
	} else if b >= a && b >= c {
		return b
	}
	return c
}
