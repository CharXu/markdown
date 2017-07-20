package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	//"math/rand"
	"net/http"
	//"strconv"
	"strings"
	"time"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"
)

type respfromsrv interface {
	Unmarshal(dAtA []byte) error
}

//服务器返回的信息
type respMsg struct {
	respbytes int64
	resptime  float64
	body      respfromsrv
}

func main() {

	//命令行参数
	cFlag := flag.Int("c", 1, "并发的连接数concurrent connects")
	nFlag := flag.Int("n", 1, "请求次数")
	// tFlag := flag.Int("t", 0, "请求间隔时间(毫秒)")
	uFlag := flag.String("u", "http://192.168.0.51/char", "测试的URL，格式：http://hostname")
	mFlag := flag.String("m", "POST", "http的请求方法")
	rFlag := flag.String("r", "hello", "which request you want to test")

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
	var serverresp respfromsrv

	ch := make(chan *respMsg, *cFlag)
	eachRequest := *nFlag / (*cFlag)

	start := time.Now()
	for i := 0; i < *cFlag; i++ {
		go GoRequest(*uFlag, *mFlag, eachRequest, *rFlag, ch)
	}

	for i := 0; i < *cFlag; i++ {
		result := <-ch
		//todo
		totalRequestTime += result.resptime
		totalTransferBytes += result.respbytes
		serverresp = result.body
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

// GoRequest ...
// 发送请求
func GoRequest(url string, method string, eachRequest int, reqtype string, ch chan *respMsg) {
	client := &http.Client{}
	respmsg := &respMsg{}

	var reqreader io.Reader
	//var firstreader io.Reader
	// var urlfirst string
	//var req http.Request

	switch reqtype {
	case "hello":
		helloreader, err := sethelloPkt()
		if err != nil {
			fmt.Println("hello pakect err: ", err)
			return
		}
		reqreader = helloreader
		url += "/webHello"

	case "login":
		loginreader, err := setloginPkt()
		if err != nil {
			fmt.Println("login packet err: ", err)
			return
		}
		reqreader = loginreader
		url += "/webLogin"

	case "build":
		buildreader, err := setbuildPkt()
		if err != nil {
			fmt.Println("build packet err: ", err)
			return
		}
		// loginreader, err := setloginPkt()
		// if err != nil {
		// 	fmt.Println("login packet err :", err)
		// 	return
		// }
		// // firstreader = loginreader
		// urlfirst = url + "/webLogin"
		reqreader = buildreader
		url += "/webRequest"
	}

	var timeonerequest float64
	var bytesonerequest int64
	//var bodyonerequest []byte
	var respPkt afproto.Packet
	var respsrv respfromsrv

	for i := 0; i < eachRequest; i++ {
		req, err := http.NewRequest(method, url, reqreader)
		if err != nil {
			fmt.Println("set request set err: ", err)
			return
		}
		req.Header.Set("Connection", "Keep-Alive")

		// if reqtype == "build" {
		// 	reqfirst, err := http.NewRequest(method, urlfirst, firstreader)
		// 	if err != nil {
		// 		fmt.Println("set first request err: ", err)
		// 		return
		// 	}
		// 	//reqfirst.Header.Set("Connection", "Keep-Alive")
		// 	respfirst, err := client.Do(reqfirst)
		// 	if err != nil {
		// 		fmt.Println("first response err: ", err)
		// 		return
		// 	}
		// 	respfirst.Body.Close()
		// }

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
		}

		respbytes, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			fmt.Println("io.Copy err: ", err)
		}
		resp.Body.Close()
		timeonerequest += resptime
		bytesonerequest += respbytes

		respdecode, _, err := cpt.Base64Decode(body)
		if err != nil {
			fmt.Println("解码失败", err)
			goto out
		}

		err = respPkt.Unmarshal(respdecode)
		if err != nil {
			fmt.Println("Pkt unmarshal err: ", err)
			goto out
		}

		respPktbody := respPkt.Body

		switch reqtype {
		case "hello":
			var respbody tdproto.HelloRsp
			err = respbody.Unmarshal(respPktbody)
			if err != nil {
				fmt.Println("rsphellobody unmarshal err: ", err)
				goto out
			}
			respsrv = &respbody
		case "login":
			var respbody tdproto.LoginRsp
			err = respbody.Unmarshal(respPktbody)
			if err != nil {
				fmt.Println("rsploginbody unmarshal err: ", err)
				goto out
			}
			respsrv = &respbody
		case "build":
			var respbody tdproto.BuildRsp
			err = respbody.Unmarshal(respPktbody)
			if err != nil {
				fmt.Println("rspbuildbody unmarshal err: ", err)
				goto out
			}
			respsrv = &respbody
		}

		// time.Sleep(10 * time.Millisecond)
	}

out:
	respmsg.body, respmsg.resptime, respmsg.respbytes = respsrv, timeonerequest, bytesonerequest
	ch <- respmsg

}

//setHelloPkt ...
//hello请求
func sethelloPkt() (io.Reader, error) {
	helloreq := &tdproto.HelloReq{}
	hellobody, err := helloreq.Marshal()
	if err != nil {
		return nil, err
	}

	helloPkt := &afproto.Packet{
		Head: &afproto.PktHead{
			Opt: &afproto.PktOpt{
				Version: 1,
				Mtkey:   1,
				Skey:    1,
				Seq:     1,
				Time:    1,
				UID:     "121212",
				Device:  "1",
				Channel: "1",
				Lang:    "Zh",
			},
			PktInfo: [](*afproto.PktInfo){
				&afproto.PktInfo{
					Cmd: uint32(tdproto.GAME_CMD_HELLO_REQ),
					Len: uint32(len(hellobody)),
				},
			},
		},
		Body: hellobody,
	}
	hellobytes, err := helloPkt.Marshal()
	if err != nil {
		fmt.Println("Hello packect marshal err: ", err)
		return nil, err
	}

	hellobytescpt, err := cpt.Base64Encode(hellobytes)
	if err != nil {
		return nil, err
	}

	helloreader := bytes.NewReader(hellobytescpt)

	return helloreader, nil
}

type Pkt interface {
	Marshal() (dAtA []byte, err error)
}

//setlPkt ...
// func setPkt(reqtype string) (io.Reader, error) {
// 	var req
// 	switch reqtype {
// 		case "hello":
// 			helloreq := &tdproto.HelloReq{}
// 	}
// }

func setloginPkt() (io.Reader, error) {
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// token := strconv.Itoa(r.Intn(100000000))
	loginreq := &tdproto.LoginReq{
		Type:        tdproto.LOGIN_TYPE_GUEST,
		AccessToken: "2222dsfa",
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
	loginbytes, err := loginPkt.Marshal()
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

func setbuildPkt() (io.Reader, error) {
	buildreq := &tdproto.BuildReq{
		Index:   2,
		ToLevel: 1,
	}
	buildbody, err := buildreq.Marshal()
	if err != nil {
		return nil, err
	}

	buildPkt := &afproto.Packet{
		Head: &afproto.PktHead{
			Opt: &afproto.PktOpt{
				Version: 1,
				Mtkey:   1,
				Skey:    1,
				Seq:     1,
				Time:    1,
				UID:     "100031",
				Device:  "1",
				Channel: "1",
				Lang:    "Zh",
			},
			PktInfo: [](*afproto.PktInfo){
				&afproto.PktInfo{
					Cmd: uint32(tdproto.GAME_CMD_BUILD_REQ),
					Len: uint32(len(buildbody)),
				},
			},
		},
		Body: buildbody,
	}
	buildbytes, err := buildPkt.Marshal()
	if err != nil {
		return nil, err
	}

	buildbytescpt, err := cpt.Base64Encode(buildbytes)
	if err != nil {
		return nil, err
	}

	buildreader := bytes.NewReader(buildbytescpt)

	return buildreader, nil
}
