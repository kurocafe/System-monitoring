# Q6. Slice + Function

Implement this function:
```go
func Filter(nums []int, fn func(int) bool) []int
```

Usage:
```go
// keep only even numbers
Filter([]int{1,2,3,4}, func(n int) bool { return n%2==0 })
```

Expected output: `[2,4]`