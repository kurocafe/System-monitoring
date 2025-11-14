package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    fmt.Println(add(2,3))
    fmt.Println(naked_return(4,3))
    for i := 0; i < 10; i ++ {
        // main がreturn をしたタイミングでdefer してた関数が新しい順で実行される？
        defer fmt.Printf("nachinchin %d \n", i)
    }

    fmt.Println(pow(3, 1, 5))
    x := float64(5)
    sqrt(x)
}