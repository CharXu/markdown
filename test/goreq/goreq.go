package goreq

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"char/markdown/test/getpkt"
	"char/markdown/test/islandBuild"
	"char/markdown/test/mystruct"
	"char/markdown/test/setpkt"

	"sync"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"
)

// GoRequest ...
// 发送请求
func GoRequest(url string, method string, eachRequest int, reqtype string, ch chan *mystruct.RespMsg, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	respmsg := &mystruct.RespMsg{}

	var reqbytes []byte
	switch reqtype {
	case "hello":
		hellocpt, err := setpkt.SethelloPkt()
		if err != nil {
			fmt.Println("hello pakect err: ", err)
			return
		}
		reqbytes = hellocpt
		url += "/webHello"

	case "login":
		logincpt, err := setpkt.SetloginPkt()
		if err != nil {
			fmt.Println("login packet err: ", err)
			return
		}
		reqbytes = logincpt
		url += "/webLogin"

	default:
		buildcpt, err := islandBuild.SetbuildPkt(method, url, 0, 1, true)
		if err != nil {
			fmt.Println("build packet err: ", err)
			return
		}
		reqbytes = buildcpt
		url += "/webRequest"
	}

	var timeonerequest float64
	var bytesonerequest int64
	//var bodyonerequest []byte
	//var respPkt afproto.Packet
	var respsrv mystruct.Respfromsrv

	for i := 0; i < eachRequest; i++ {
		reqreader := bytes.NewReader(reqbytes)
		req, err := http.NewRequest(method, url, reqreader)
		if err != nil {
			fmt.Println("set request set err: ", err)
			return
		}
		req.Header.Set("Connection", "Keep-Alive")

		//开始请求
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("get response err: ", err)
			return
		}
		resptime := time.Since(start).Seconds()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("response Read err: ", err)
			return
		}

		respbytes, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			fmt.Println("io.Copy err: ", err)
			return
		}
		resp.Body.Close()
		timeonerequest += resptime
		bytesonerequest += respbytes

		rspbody, err := getpkt.ByteToPkt(reqtype, body)
		if err != nil {
			fmt.Println("byte decode err: ", err)
			return
		}
		respsrv = rspbody

		if value, ok := respsrv.(*tdproto.BuildRsp); ok {

			if value.Island == nil {
				fmt.Println("read island info fail, no info get")
				return
			}

			buildresp := new(tdproto.BuildRsp)
			buildresp = value
			list := buildresp.Island.GetBDList()
			reqbytes, err = islandBuild.Islandlogic(method, url, list)
			if err != nil {
				fmt.Println("build logic err : ", err)
				return
			}

		}
		respmsg.Body, respmsg.Resptime, respmsg.Respbytes = respsrv, timeonerequest, bytesonerequest
		ch <- respmsg

		// time.Sleep(10 * time.Millisecond)
	}
}
