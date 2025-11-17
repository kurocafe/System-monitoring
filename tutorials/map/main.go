package main

import "fmt"

func main(){
	words := []string{"apple", "banana", "apple", "orange", "banana", "apple"}

	num_words := make(map[string]int)

	for _, word := range words {
		element, exist := num_words[word]
		if exist {
			num_words[word] = element + 1
		}else{
			num_words[word] = 1
		}
	}

	fmt.Println(num_words)
}