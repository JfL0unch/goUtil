package utils

type Number interface {
	int |
		int8 |
		int16 |
		int32 |
		int64 |
		uint |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		float32 |
		float64
}

func MaxVal[T Number](inputs []T) T {
	if len(inputs) == 0 {
		return 0
	}
	max := inputs[0]
	for _, v := range inputs {
		if v > max {
			max = v
		}
	}
	return max
}
func SumVal[T Number](inputs []T) T {
	if len(inputs) == 0 {
		return 0
	}
	sum := inputs[0]
	sum = 0

	for _, v := range inputs {
		sum += v
	}
	return sum
}
func MinVal[T Number](inputs []T) T {
	if len(inputs) == 0 {
		return 0
	}
	min := inputs[0]
	for _, v := range inputs {
		if v < min {
			min = v
		}
	}
	return min
}

// LikeEqual 近似相等
// ratio 相似范围
func LikeEqual[T Number](x1, x2, delta T) bool {
	if x1 < x2 {
		return x2-x1 < delta
	} else {
		return x1-x2 < delta
	}
}
