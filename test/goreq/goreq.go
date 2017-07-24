package goreq

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/getpkt"
	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/islandBuild"
	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/mystruct"
	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/setpkt"

	"sync"

	"os"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"
)

// GoRequest ...
// 发送请求
func GoRequest(url string, method string, eachRequest int, reqtype string, ch chan *mystruct.RespMsg, wg *sync.WaitGroup, timeDur int) {
	defer wg.Done()
	client := &http.Client{}
	respmsg := &mystruct.RespMsg{}

	var reqbytes []byte
	var firstuid string
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
		buildcpt, uid, err := islandBuild.SetbuildPkt(method, url, 0, 1, "2", true)
		if err != nil {
			fmt.Println("build packet err: ", err)
			return
		}
		reqbytes = buildcpt
		url += "/webRequest"
		firstuid = uid
	}

	var timeonerequest float64
	var bytesonerequest int64
	//var bodyonerequest []byte
	//var respPkt afproto.Packet
	var respsrv mystruct.Respfromsrv

	reqreader := bytes.NewReader(reqbytes)

	for i := 0; i < eachRequest; i++ {

		reqreader.Reset(reqbytes)

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
			reqbytes, err = islandBuild.Islandlogic(method, url, firstuid, list)
			if err != nil {
				fmt.Println("build logic err : ", err)
				return
			}
		}
		time.Sleep(time.Duration(int64(timeDur)) * time.Millisecond)
	}
	respmsg.Body, respmsg.Resptime, respmsg.Respbytes = respsrv, timeonerequest, bytesonerequest
	ch <- respmsg
}

func checkErrr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal ero: %s", err.Error())
		os.Exit(1)
	}
}
