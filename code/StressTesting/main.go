package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//服务器返回的信息
type respMsg struct {
	respbytes  int64
	serverInfo string
	resptime   float64
	doclen     int64
}

func main() {

	//命令行参数
	cFlag := flag.Int("c", 100, "并发的连接数concurrent connects")
	nFlag := flag.Int("n", 10000, "请求次数")
	tFlag := flag.Int("t", 0, "请求间隔时间(毫秒)")
	uFlag := flag.String("u", "http://localhost", "测试的URL，格式：http://hostname")
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
	if *tFlag < 0 || *tFlag > 1000 {
		fmt.Println("请求间隔时间为0-1000毫秒")
		return
	}
	*mFlag = strings.ToUpper(*mFlag)
	if *mFlag != "GET" && *mFlag != "POST" {
		fmt.Println("invalid http method, only get/post supported")
		return
	}
	var respmsg respMsg

	var totalTime float64
	var totalTransferBytes int64
	var totalRequestTime int
	var client = new(http.Client)

	//并发请求
	start := time.Now()

	ch := make(chan *respMsg, *cFlag)
	eachRequest := *nFlag / (*cFlag)
	for i := 0; i < *cFlag; i++ {
		go GoRequest(client, *uFlag, *mFlag, eachRequest, ch)
	}

	for i := 0; i < *cFlag; i++ {
		result := <-ch
		//todo
	}

	totaltime := time.Since(start).Seconds()

	fmt.Println("Server Host: ", strings.TrimPrefix(*uFlag, "http://"))
	fmt.Println("Concurrent Level: ", *cFlag)
	fmt.Printf("Total Test Time: %.4fs\n", totaltime) //请求时间
	fmt.Println("timeslen: ", totalRequestTime)       //请求次数
	fmt.Println("Document length: ", totalTransferBytes/totalRequestTime)

	//平均请求时间
	var avetime float64
	for _, i := range times {
		avetime = avetime + i
	}
	fmt.Printf("Average Response Time: %.4fs\n", avetime/float64(len(times)))
	fmt.Printf("Total Response Time: %.4fs\n", avetime)

	//请求总字节数
	var totalbytes int64
	for _, i := range bytes {
		totalbytes = totalbytes + i
	}
	fmt.Println("total transferred: ", totalbytes, "bytes")
	fmt.Println("Average transferred: ", totalbytes/int64(*cFlag), "bytes")

}

// GoRequest ...
// 发送请求
func GoRequest(client *http.Client, url string, method string, eachRequest int, ch chan *respMsg) {
	respMsg := &respMsg{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("request err : ", err)
		return
	}
	// req.Header.Set("")
	// req.Header.Set()
	// req.Header.Set()
	for i := 0; i < eachRequest; i++ {
		start := time.Now()
		respServer, err := client.Do(req)
		if err != nil {
			fmt.Println("response err: ", err)
			return
		}
		resptime := time.Since(start).Seconds()
		if respServer == nil {
			return
			fmt.Println("Maybe the server is not avalible")
		}

		respbytes, err := io.Copy(ioutil.Discard, respServer.Body)
		serverInfo := respServer.Header.Get("Server")
		doclen := respServer.ContentLength
		respMsg.doclen, respMsg.serverInfo, respMsg.resptime, respMsg.respbytes = doclen, serverInfo, resptime, doclen
		respServer.Body.Close()

		// time.Sleep(10 * time.Millisecond)
	}
	ch <- respMsg

}
