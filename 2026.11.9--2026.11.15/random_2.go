package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Thing struct {
	ID  int
	NUM int
}

func main() {
	ch := make(chan Thing, 20)
	var wg sync.WaitGroup
	var c Thing

	for i := 1; i <= 20; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			sleeptime := rand.Intn(1001)
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)

			ch <- Thing{
				ID:  id,
				NUM: rand.Intn(101),
			}
		}(i)
	}
	wg.Wait()

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
