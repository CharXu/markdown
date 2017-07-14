package main

import (
	"net/http"
	"os"
	"strconv"
	"fmt"
	"time"
)

func main() {

	url := os.Args[1]
	times, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Println("string convert to int err: ", err)
	}

	//并发请求
	for i := 0; i < times; i++ {
		go GoRequest()

		PrintResp(<-ch)
	}
}

// GoRequest ...
// 发送请求
func GoRequest(url string) {
	start := time.Now()

	client := &http.Client{}

	req, err := http.NewRequest("GET", url)
	if err != nil {
		fmt.Println("request err: ", err)
		return
	}

	// req.Header.Set()
	// req.Header.Set()
	// req.Header.Set()

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("response err: ", err)
		return
	}

	resp.
	defer resp.Body.Close()

	return
}

// PrintResp ...
// 响应结果的处理
func PrintResp() {

}
