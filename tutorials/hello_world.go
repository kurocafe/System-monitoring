package main

import "fmt"

func sumOfSquares(nums []int) int {
    res := 0
    for _, n :=  range nums {
        res += n * n
    }

    return res
}


func main() {
    numbers := []int{3,4,5,6,7,8}
    sum := sumOfSquares(numbers)
    fmt.Println(sum)
}


