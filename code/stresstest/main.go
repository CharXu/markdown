package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

import tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

import afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"

import cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"

//服务器返回的信息
type respMsg struct {
	respbytes int64
	resptime  float64
	body      []byte
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
	var respbody []byte

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
		respbody = result.body
	}

	totaltime := time.Since(start).Seconds()

	fmt.Println("Server Host: ", strings.TrimPrefix(*uFlag, "http://"))
	fmt.Println("Concurrent Level: ", *cFlag)
	fmt.Printf("Total Test Time: %.4fs\n", totaltime)
	fmt.Printf("time per connection : %.4fs\n", totalRequestTime/float64(*cFlag))
	fmt.Printf("Time per request : %.4fs\n", totalRequestTime/float64(*nFlag))
	//fmt.Println("Document length: ", totalTransferBytes/int64(totalRequestTime))
	fmt.Println("response body: ", respbody)
}

// GoRequest ...
// 发送请求
func GoRequest(url string, method string, eachRequest int, ch chan *respMsg) {
	client := &http.Client{}
	respmsg := &respMsg{}

	helloreader, err := sethelloPkt()
	if err != nil {
		fmt.Println("hello pakect err: ", err)
		return
	}

	loginreader, err := setloginPkt()
	if err != nil {
		fmt.Println("login packet err: ", err)
		return
	}

	var timeonerequest float64
	var bytesonerequest int64
	var bodyonerequest []byte

	for i := 0; i < eachRequest; i++ {
		reqhello, err := http.NewRequest(method, url, helloreader)
		if err != nil {
			fmt.Println("hello request set err: ", err)
			return
		}
		reqhello.Header.Set("Connection", "Keep-Alive")

		reqlogin, err := http.NewRequest(method, url, loginreader)
		if err != nil {
			fmt.Println("login request set err : ", err)
			return
		}
		reqlogin.Header.Set("Connection", "Keep-Alive")

		start := time.Now()

		resphello, err := client.Do(reqhello)
		if err != nil {
			fmt.Println("hello response err: ", err)
			return
		}

		resplogin, err := client.Do(reqlogin)
		if err != nil {
			fmt.Println("login response err: ", err)
			return
		}

		resptime := time.Since(start).Seconds()

		resphello.Body.Close()

		// loginRespbody, err := ioutil.ReadAll(resplogin.Body)
		if err != nil {
			fmt.Println("read response err: ", err)
		}
		// loginRespPkt := &tdproto.LoginRsp{}
		// err = proto.Unmarshal(loginRespbody, loginRespPkt)
		// if err != nil {
		// 	fmt.Println("response unmarshal err: ", err)
		// 	return
		// }
		// if loginRespPkt.EncodeKey == "" {
		// 	fmt.Println("login failed")
		// }

		body, err := ioutil.ReadAll(resplogin.Body)
		if err != nil {
			fmt.Println("response Read err: ", err)
		}

		respbytes, err := io.Copy(ioutil.Discard, resplogin.Body)
		if err != nil {
			fmt.Println("io.Copy err: ", err)
		}
		resplogin.Body.Close()

		timeonerequest += resptime
		bytesonerequest += respbytes
		bodyonerequest = body

		// time.Sleep(10 * time.Millisecond)
	}

	respmsg.body, respmsg.resptime, respmsg.respbytes = bodyonerequest, timeonerequest, bytesonerequest
	ch <- respmsg

}

//setHelloPkt ...
//hello请求
func sethelloPkt() (io.Reader, error) {

	helloPkt := &tdproto.HelloReq{}
	hellobytes, err := helloPkt.Marshal()
	if err != nil {
		fmt.Println("Hello packect marshal err: ", err)
		return nil, err
	}

	// hellobytescpt, err := cpt.Base64Encode(hellobytes)
	// if err != nil {
	// 	return nil, err
	// }

	helloreader := bytes.NewReader(hellobytes)

	return helloreader, nil
}

//setloginPkt ...
//login request
func setloginPkt() (io.Reader, error) {
	loginreq := &tdproto.LoginReq{
		Type:        tdproto.LOGIN_TYPE_GUEST,
		AccessToken: "2",
		Gender:      tdproto.GENDER_TYPE_MALE,
	}
	loginbody, err := loginreq.Marshal()
	if err != nil {
		return nil, err
	}

	loginPkt := &afproto.Packet{
		Head: &afproto.PktHead{
			Opt: &afproto.PktOpt{
				Version: 1,
				Mtkey:   1,
				Skey:    1,
				Seq:     1,
				Time:    1,
				UID:     "",
				Device:  "1",
				Channel: "1",
				Lang:    "Zh",
			},
			PktInfo: [](*afproto.PktInfo){
				&afproto.PktInfo{
					Cmd: uint32(tdproto.GAME_CMD_LOGIN_REQ),
					Len: uint32(len(loginbody)),
				},
			},
		},
		Body: loginbody,
	}
	loginbytes, err := proto.Marshal(loginPkt)
	if err != nil {
		return nil, err
	}

	loginbytescpt, err := cpt.Base64Encode(loginbytes)
	if err != nil {
		return nil, err
	}

	loginreader := bytes.NewReader(loginbytescpt)

	return loginreader, nil
}
