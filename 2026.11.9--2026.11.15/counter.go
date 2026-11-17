package main

import "fmt"

func main() {
	count1 := Counter()
	var sum int

	for i := 0; i < 5; i++ {
		sum = count1()
	}

	fmt.Printf("The result is %d.", sum)
}

func Counter() func() int {
	count := 0
	return func() int {
		count += 1
		return count
	}
}
