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
	resptime   float64
	serverInfo string
	doclen     int64
}

func main() {

	//命令行参数
	cFlag := flag.Int("c", 1, "并发的连接数concurrent connects")
	uFlag := flag.String("u", "http://localhost", "测试的URL")
	mFlag := flag.String("m", "GEt", "http的请求方法")

	fmt.Println("method", *mFlag)

	flag.Parse()

	var respmsg respMsg

	var times []float64
	var bytes []int64

	//并发请求
	start := time.Now()

	ch := make(chan respMsg)

	for i := 0; i < *cFlag; i++ {
		go GoRequest(*uFlag, ch)
		respmsg = <-ch
		times = append(times, respmsg.resptime)
		bytes = append(bytes, respmsg.respbytes)
	}

	totaltime := time.Since(start).Seconds()

	fmt.Println("ServerInfo: ", respmsg.serverInfo)
	fmt.Println("Server Host: ", strings.TrimPrefix(*uFlag, "http://"))
	fmt.Println("Concurrent Level: ", *cFlag)
	fmt.Printf("Total Time: %.2fs\n", totaltime) //请求时间
	fmt.Println("timeslen: ", len(times))        //请求次数
	fmt.Println("Document length: ", respmsg.doclen)

	//平均请求时间
	var avetime float64
	for _, i := range times {
		avetime = avetime + i
	}
	fmt.Printf("averagetime: %.4fs\n", avetime/float64(len(times)))

	//请求总字节数
	var totalbytes int64
	for _, i := range bytes {
		totalbytes = totalbytes + i
	}
	fmt.Println("total transferred: ", totalbytes)

}

// GoRequest ...
// 发送请求
func GoRequest(url string, ch chan respMsg) {
	//var reqs *http.Request

	start := time.Now()

	// client := &http.Client{}
	// b := strings.NewReader("name=cjb")

	// switch method {
	// case "GET":
	// 	req, err := http.NewRequest(method, url, b)
	// 	if err != nil {
	// 		fmt.Println("request err : ", err)
	// 	}
	// 	reqs = req
	// }

	// req.Header.Set()
	// req.Header.Set()
	// req.Header.Set()

	respServer, err := http.Get(url)
	defer respServer.Body.Close()
	if err != nil {
		fmt.Println("response err: ", err)
		return
	}
	respbytes, err := io.Copy(ioutil.Discard, respServer.Body)

	serverInfo := respServer.Header.Get("Server")

	doclen := respServer.ContentLength

	resptime := time.Since(start).Seconds()
	respmsg := respMsg{respbytes, resptime, serverInfo, doclen}

	ch <- respmsg
}

// DealResp ...
// 响应结果的处理
// func DealResp(resp respMsg) {
// 	respbytes, err := io.Copy(ioutil.Discard, resp.respServer.Body)
// 	if err != nil {
// 		fmt.Println("io.Copy err: ", err)
// 	}
// 	fmt.Println(respbytes)
// }
