package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func request(i int) {

	reqData := fmt.Sprintf("{\"seq\": %d}", i)

	req, _ := http.NewRequest("GET", "http://117.72.42.214:8080/api/reward", bytes.NewBuffer([]byte(reqData)))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdGFnZSI6InJld2FyZCIsInNpZCI6IjhhNmQxYjA4ZjI4YTQ2NTc3Y2ZkNDUwYTJlZTIzOTU2IiwiZXhwIjoxNzk1MTkyMzM5LCJpYXQiOjE3OTUxOTE3MzksImp0aSI6ImQ3NWNiMTdmMzhmYzQxNDE5NDkzMWVmZjkxZTg1YmU1In0.iu4Ka3wkGj_7XKGjUouWwiCePGu7twOEkAummhipjz8")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("收到回复:", string(body))
}

func main() {
	nums := []int{1, 5, 3, 2, 4, 7, 8, 9, 10, 6}
	var wg sync.WaitGroup

	for _, v := range nums {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			request(id)
		}(v)
	}
	wg.Wait()

}
