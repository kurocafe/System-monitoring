package main

import (
	"fmt"
	"math/rand"
)

func sqrt(x float64) float64{
	z := rand.Float64()
	fmt.Printf("initialized: %g \n", z)

	for i := 0; i < 10; i ++ {
		z -= (z*z - x) / (2*z)
		fmt.Printf("%d: %g \n", i + 1, z)
	}

	return z
}