package goreq

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/setpkt"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"

	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"
)

type Respfromsrv interface {
	Unmarshal(dAtA []byte) error
}

//服务器返回的信息
type RespMsg struct {
	Respbytes int64
	Resptime  float64
	Body      Respfromsrv
}


// GoRequest ...
// 发送请求
func GoRequest(url string, method string, eachRequest int, reqtype string, ch chan *RespMsg) {
	client := &http.Client{}
	respmsg := &RespMsg{}

	var reqreader io.Reader
	//var firstreader io.Reader
	// var urlfirst string
	//var req http.Request

	switch reqtype {
	case "hello":
		helloreader, err := setpkt.SethelloPkt()
		if err != nil {
			fmt.Println("hello pakect err: ", err)
			return
		}
		reqreader = helloreader
		url += "/webHello"

	case "login":
		loginreader, err := setpkt.SetloginPkt()
		if err != nil {
			fmt.Println("login packet err: ", err)
			return
		}
		reqreader = loginreader
		url += "/webLogin"

	case "build":
		buildreader, err := setpkt.SetbuildPkt()
		if err != nil {
			fmt.Println("build packet err: ", err)
			return
		}
		// loginreader, err := setpkt.SetloginPkt()
		// if err != nil {
		// 	fmt.Println("login packet err :", err)
		// 	return
		// }
		// firstreader = loginreader
		// urlfirst = url + "/webLogin"
		reqreader = buildreader
		url += "/webRequest"
	}

	var timeonerequest float64
	var bytesonerequest int64
	//var bodyonerequest []byte
	//var respPkt afproto.Packet
	var respsrv Respfromsrv

	for i := 0; i < eachRequest; i++ {
		req, err := http.NewRequest(method, url, reqreader)
		if err != nil {
			fmt.Println("set request set err: ", err)
			return
		}
		req.Header.Set("Connection", "Keep-Alive")

		// if reqtype == "build" {
		// 	loginfirst, err := http.NewRequest(method, urlfirst, firstreader)
		// 	if err != nil {
		// 		fmt.Println("set first request err: ", err)
		// 		return
		// 	}
		// 	//reqfirst.Header.Set("Connection", "Keep-Alive")
		// 	respfirst, err := client.Do(loginfirst)
		// 	if err != nil {
		// 		fmt.Println("first response err: ", err)
		// 		return
		// 	}
		// 	bodyisland, err := ioutil.ReadAll(respfirst.Body)
		// 	if err != nil {
		// 		fmt.Println("first login rsp err: ", err)
		// 		return
		// 	}
		// 	respfirst.Body.Close()
		// 	islanddecode, _, err := cpt.Base64Decode(bodyisland)
		// 	if err != nil {
		// 		fmt.Println("first decode err: ", err)
		// 		return
		// 	}

		// 	var islandPkt afproto.Packet
		// 	err = islandPkt.Unmarshal(islanddecode)
		// 	if err != nil {
		// 		fmt.Println("Island Pkt unmarshal err: ", err)
		// 		return
		// 	}
		// 	islandPktbody := islandPkt.Body
		// 	var islandbody

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

		rspbody, err := ByteToPkt(reqtype, body)
		if err != nil {
			fmt.Println("byte decode err: ", err)
		}

		switch reqtype {
			case "hello":
			var hellorspbody *tdproto.HelloRsp
			hellorspbody = rspbody
			respsrv = hellorspbody
		}

		// time.Sleep(10 * time.Millisecond)
	}

	//out:
	respmsg.Body, respmsg.Resptime, respmsg.Respbytes = , timeonerequest, bytesonerequest
	ch <- respmsg
}

// ByteToPkt ...
//还原proto结构的数据
func ByteToPkt(reqtype string, rspbytes []byte) (Respfromsrv, error) {
	respdecode, _, err := cpt.Base64Decode(rspbytes)
	if err != nil {
		return nil, err
	}

	var rspPkt *afproto.Packet

	err = rspPkt.Unmarshal(respdecode)
	if err != nil {
		return nil, err
	}
	rspPktbody := rspPkt.Body

	switch reqtype {
	case "hello":
		var rspbody tdproto.HelloRsp
		err = rspbody.Unmarshal(rspPktbody)
		if err != nil {
			return nil, err
		}
		return &rspbody, nil
	case "login":
		var rspbody tdproto.LoginRsp
		err = rspbody.Unmarshal(rspPktbody)
		if err != nil {
			return nil, err
		}
		return &rspbody, nil
	case "build":
		var rspbody tdproto.BuildRsp
		err = rspbody.Unmarshal(rspPktbody)
		if err != nil {
			return nil, err
		}
		return &rspbody, nil
	default:
		return nil, fmt.Errorf("invalid reqtype")
	}
}
