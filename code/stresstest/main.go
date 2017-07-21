package main

import (
	"flag"
	"fmt"
	//"math/rand"

	//"strconv"
	"strings"
	"time"

	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/goreq"
)

func main() {

	//命令行参数
	cFlag := flag.Int("c", 2, "并发的连接数concurrent connects")
	nFlag := flag.Int("n", 4, "请求次数")
	// tFlag := flag.Int("t", 0, "请求间隔时间(毫秒)")
	uFlag := flag.String("u", "http://192.168.0.51/char", "测试的URL，格式：http://hostname")
	mFlag := flag.String("m", "POST", "http的请求方法")
	rFlag := flag.String("r", "build", "which request you want to test")

	flag.Parse()

	if *cFlag >= 10000 || *cFlag < 0 {
		fmt.Println("param connections invalid, should be between 1 and 9999 ")
		return
	}
	if *nFlag < 0 || *nFlag > 10000000 {
		fmt.Println("请求次数范围0-10000000")
		return
	}

	// if *tFlag < 0 || *tFlag > 1000 {
	// 	fmt.Println("请求间隔时间为0-1000毫秒")
	// 	return
	// }

	*mFlag = strings.ToUpper(*mFlag)
	if *mFlag != "GET" && *mFlag != "POST" {
		fmt.Println("invalid http method, only get/post supported")
		return
	}

	*rFlag = strings.ToLower(*rFlag)
	if *rFlag != "hello" && *rFlag != "login" && *rFlag != "build" {
		fmt.Println("invalid request, only hello/login/build supported")
		return
	}

	var totalTransferBytes int64
	var totalRequestTime float64
	var serverresp goreq.Respfromsrv

	ch := make(chan *goreq.RespMsg, *cFlag)
	eachRequest := *nFlag / (*cFlag)

	start := time.Now()
	for i := 0; i < *cFlag; i++ {
		go goreq.GoRequest(*uFlag, *mFlag, eachRequest, *rFlag, ch)
	}

	for i := 0; i < *cFlag; i++ {
		result := <-ch
		//todo
		totalRequestTime += result.Resptime
		totalTransferBytes += result.Respbytes
		serverresp = result.Body
	}

	totaltime := time.Since(start).Seconds()

	fmt.Println("Server Host: ", strings.TrimPrefix(*uFlag, "http://"))
	fmt.Println("Concurrent Level: ", *cFlag)
	fmt.Printf("Total Test Time: %.4fs\n", totaltime)
	fmt.Printf("time per connection : %.4fs\n", totalRequestTime/float64(*cFlag))
	fmt.Printf("Time per request : %.4fs\n", totalRequestTime/float64(*nFlag))
	//fmt.Println("Document length: ", totalTransferBytes/int64(*nFlag))
	fmt.Printf("response body: %v\n", serverresp)
}
