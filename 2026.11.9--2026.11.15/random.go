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
	close(ch)

	var arr []Thing
	for thing := range ch {
		arr = append(arr, thing)
	}
	fmt.Println("排序前", arr)
	Sort(arr)
}

func Sort(arr []Thing) {
	for i := 0; i < 20; i++ {
		for arr[i].ID != i+1 {
			arr[i], arr[arr[i].ID-1] = arr[arr[i].ID-1], arr[i]
		}
	}
	fmt.Println("排序后", arr)
}
