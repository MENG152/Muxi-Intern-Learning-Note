package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Thing struct {
	ID  int
	NUM int
}

func main() {
	ch := make(chan Thing, 20)
	var c Thing
	doneCH := make(chan struct{})

	for i := 1; i <= 20; i++ {

		go func(id int) {
			defer func() { doneCH <- struct{}{} }()

			sleeptime := rand.Intn(1001)
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)

			ch <- Thing{
				ID:  id,
				NUM: rand.Intn(101),
			}
		}(i)
	}
	for i := 0; i < 20; i++ {
		<-doneCH
	}

	fmt.Println("排序前")

	for i := 0; i < 20; i++ {
		c = <-ch
		fmt.Printf("ID is %d, NUM is %d.\n", c.ID, c.NUM)
		ch <- c
	}

	fmt.Println("排序后")

	for i := 0; i < 20; i++ {
		for j := 0; j < 20-i; j++ {
			t := <-ch
			if t.ID == i+1 {
				fmt.Printf("ID is %d, NUM is %d.\n", t.ID, t.NUM)
				break
			}
			ch <- t
		}
	}
	close(ch)
}
