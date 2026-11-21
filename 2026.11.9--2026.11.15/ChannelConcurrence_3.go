package main

import (
	"fmt"
)

func main() {
	s1 := "0123456789"
	s2 := "ABCDEFGHIJKLMN"

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	endCH := make(chan struct{})

	go func() {
		defer func() {
			endCH <- struct{}{}
		}()
		for i := 0; i < len(s1); i += 2 {
			_, ok := <-ch1
			if !ok {
				fmt.Print(s1[i:])
				break
			}
			end := min(i+2, len(s1))
			fmt.Print(s1[i:end])
			if end == len(s1) {
				close(ch2)
				break
			} else {
				ch2 <- struct{}{}
			}

		}
	}()

	go func() {
		defer func() {
			endCH <- struct{}{}
		}()
		for i := 0; i < len(s2); i += 2 {
			_, ok := <-ch2
			if !ok {
				fmt.Print(s2[i:])
				break
			}
			end := min(i+2, len(s2))
			fmt.Print(s2[i:end])
			if end == len(s2) {
				close(ch1)
				break
			} else {
				ch1 <- struct{}{}
			}
		}
	}()

	ch1 <- struct{}{}

	for i := 0; i < 2; i++ {
		<-endCH
	}

	close(endCH)
}
