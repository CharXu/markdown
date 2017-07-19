package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	pbtest "char/markdown/code/StressTesting/pbtest"

	"bytes"

	"github.com/golang/protobuf/proto"
)

//服务器返回的信息
type respMsg struct {
	respbytes int64
	resptime  float64
	doclen    int64
}

func main() {

	//命令行参数
	cFlag := flag.Int("c", 100, "并发的连接数concurrent connects")
	nFlag := flag.Int("n", 10000, "请求次数")
	// tFlag := flag.Int("t", 0, "请求间隔时间(毫秒)")
	uFlag := flag.String("u", "http://www.cnblog.com", "测试的URL，格式：http://hostname")
	mFlag := flag.String("m", "GET", "http的请求方法, 暂时只有GET")

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

	var totalTransferBytes int64
	var totalRequestTime float64

	ch := make(chan *respMsg, *cFlag)
	eachRequest := *nFlag / (*cFlag)

	start := time.Now()
	for i := 0; i < *cFlag; i++ {
		go GoRequest(*uFlag, *mFlag, eachRequest, ch)
	}

	for i := 0; i < *cFlag; i++ {
		result := <-ch
		//todo
		totalRequestTime += result.resptime
		totalTransferBytes += result.respbytes
	}

	totaltime := time.Since(start).Seconds()

	fmt.Println("Server Host: ", strings.TrimPrefix(*uFlag, "http://"))
	fmt.Println("Concurrent Level: ", *cFlag)
	fmt.Printf("Total Test Time: %.4fs\n", totaltime)
	fmt.Printf("time per connection : %.4fs\n", totalRequestTime/float64(*cFlag))
	fmt.Printf("Time per request : %.4fs\n", totalRequestTime/float64(*nFlag))
	//fmt.Println("Document length: ", totalTransferBytes/int64(totalRequestTime))
}

// GoRequest ...
// 发送请求
func GoRequest(url string, method string, eachRequest int, ch chan *respMsg) {
	client := &http.Client{}
	respmsg := &respMsg{}

	men := &pbtest.Men{
		Gender: "men",
		Height: 175,
		Weight: 65,
	}

	reqbytes, err := proto.Marshal(men)
	if err != nil {
		fmt.Println("marshal err : ", err)
		return
	}

	reqreader := bytes.NewReader(reqbytes)

	req, err := http.NewRequest(method, url, reqreader)
	if err != nil {
		fmt.Println("request err : ", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	// req.Header.Set()
	// req.Header.Set()

	var timeonerequest float64
	var bytesonerequest int64
	var doclenonerequest int64

	for i := 0; i < eachRequest; i++ {

		req, err := http.NewRequest(method, url, reqreader)
		if err != nil {
			fmt.Println("request err : ", err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

		start := time.Now()
		respServer, err := client.Do(req)
		if err != nil {
			fmt.Println("response err: ", err)
			return
		}
		resptime := time.Since(start).Seconds()

		respbytes, err := io.Copy(ioutil.Discard, respServer.Body)
		respServer.Body.Close()
		doclen := respServer.ContentLength

		timeonerequest += resptime
		bytesonerequest += respbytes
		doclenonerequest += doclen

		// time.Sleep(10 * time.Millisecond)
	}

	respmsg.doclen, respmsg.resptime, respmsg.respbytes = doclenonerequest, timeonerequest, bytesonerequest
	ch <- respmsg

}
