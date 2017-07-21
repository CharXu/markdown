package goreq

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/islandBuild"

	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/setpkt"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"

	"bytes"

	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"
)

type Respfromsrv interface {
	Unmarshal(dAtA []byte) error
}

//RespMsg ...
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

	var reqbytes []byte
	//var firstreader io.Reader
	// var urlfirst string
	//var req http.Request

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

	case "build":
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
	var respsrv Respfromsrv

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

		rspbody, err := ByteToPkt(reqtype, body)
		if err != nil {
			fmt.Println("byte decode err: ", err)
			return
		}
		respsrv = rspbody

		if value, ok := respsrv.(*tdproto.BuildRsp); ok {

			// if value.Island == nil {
			// 	fmt.Println("read island info fail, no info get")
			// 	return
			// }

			buildresp := new(tdproto.BuildRsp)
			buildresp = value
			list := buildresp.Island.GetBDList()
			reqbytes, err = islandBuild.Islandlogic(method, url, list)
			if err != nil {
				fmt.Println("build logic err : ", err)
			}
		}

		// time.Sleep(10 * time.Millisecond)
	}

	//out:
	respmsg.Body, respmsg.Resptime, respmsg.Respbytes = respsrv, timeonerequest, bytesonerequest
	ch <- respmsg

}

// ByteToPkt ...
//还原proto结构的数据
func ByteToPkt(reqtype string, rspbytes []byte) (Respfromsrv, error) {
	respdecode, _, err := cpt.Base64Decode(rspbytes)
	if err != nil {
		return nil, err
	}

	var rspPkt afproto.Packet

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
