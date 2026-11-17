package main

import "fmt"

type ChannleMutex chan struct{}

func (m ChannleMutex) Lock() {
	m <- struct{}{}
}

func (m ChannleMutex) Unlock() {
	<-m
}

func main() {
	var n int
	mu := make(ChannleMutex, 1)
	s1 := "123456789"
	s2 := "ABCDEFGHIJKLMN"
	doneCH := make(chan struct{})
	n = 2
	turn := 1

	go func() {
		defer func() { doneCH <- struct{}{} }()
		for i := 0; i < len(s1); i += 2 {
			for {
				mu.Lock()
				if turn == 1 || i-2 >= len(s2) {
					end := min(len(s1), i+2)
					fmt.Print(s1[i:end])
					turn = 2
					mu.Unlock()
					break
					mu.Unlock()
				}
			}
		}
	}()

	go func() {
		defer func() { doneCH <- struct{}{} }()
		for i := 0; i < len(s2); i += 2 {
			for {
				mu.Lock()
				if turn == 2 || i >= len(s1) {
					end := min(len(s2), i+2)
					fmt.Print(s2[i:end])
					turn = 1
					mu.Unlock()
					break
				}
				mu.Unlock()
			}
		}
	}()

	for i := 0; i < n; i++ {
		<-doneCH
	}
	close(doneCH)

}

func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
