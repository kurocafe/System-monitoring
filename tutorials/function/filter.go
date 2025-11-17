package main

func Filter(nums []int, fn func(int) bool) []int {
	result := make([]int, 0)

	for _, value := range nums {
		if fn(value) {
			result = append(result, value)
		}
	}

	return result
}