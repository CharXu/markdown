package islandBuild

import (
	"bytes"
	"io/ioutil"
	"net/http"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"

	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"

	"aladinfun.com/TripleDream/TripleDreamServer/tools/stresstest/setpkt"
)

//SetbuildPkt ...
//build前线登录获取UID

var uid string

func SetbuildPkt(method string, url string, no uint32, tolevel uint32, first bool) ([]byte, error) {

	if first {
		client := &http.Client{}
		logincpt, err := setpkt.SetloginPkt()
		if err != nil {
			return nil, err
		}
		url = url + "/webLogin"

		loginreader := bytes.NewReader(logincpt)
		loginreq, err := http.NewRequest(method, url, loginreader)
		if err != nil {
			return nil, err
		}

		loginresp, err := client.Do(loginreq)
		if err != nil {
			return nil, err
		}
		loginbodybytes, err := ioutil.ReadAll(loginresp.Body)
		if err != nil {
			return nil, err
		}
		loginresp.Body.Close()

		loginbodydecode, _, err := cpt.Base64Decode(loginbodybytes)
		if err != nil {
			return nil, err
		}

		var loginrspPkt afproto.Packet
		err = loginrspPkt.Unmarshal(loginbodydecode)
		if err != nil {
			return nil, err
		}
		loginPktbody := loginrspPkt.Body
		var loginbody tdproto.LoginRsp
		err = loginbody.Unmarshal(loginPktbody)
		if err != nil {
			return nil, err
		}
		uid = loginbody.UserData.UID

	}

	buildreq := &tdproto.BuildReq{
		Index:   no,
		ToLevel: tolevel,
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
				UID:     uid,
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
	return buildbytescpt, nil
}

func Islandlogic(method string, url string, islandlist []uint32) ([]byte, error) {
	// for v, _ := range islandlist {
	// 	if v < 0 || v > 4 {
	// 		return nil, fmt.Errorf("islandlist err")
	// 	}
	// }

	var reqbytes []byte
	if islandlist[4] == 5 {
		buildcpt, err := SetbuildPkt(method, url, 0, 1, false)
		if err != nil {
			return nil, err
		}
		reqbytes = buildcpt
	} else if islandlist[0] < 5 {
		buildcpt, err := SetbuildPkt(method, url, 0, islandlist[0]+1, false)
		if err != nil {
			return nil, err
		}
		reqbytes = buildcpt
	} else if islandlist[1] < 5 {
		buildcpt, err := SetbuildPkt(method, url, 1, islandlist[1]+1, false)
		if err != nil {
			return nil, err
		}
		reqbytes = buildcpt
	} else if islandlist[2] < 5 {
		buildcpt, err := SetbuildPkt(method, url, 2, islandlist[2]+1, false)
		if err != nil {
			return nil, err
		}
		reqbytes = buildcpt
	} else if islandlist[3] < 5 {
		buildcpt, err := SetbuildPkt(method, url, 3, islandlist[3]+1, false)
		if err != nil {
			return nil, err
		}
		reqbytes = buildcpt
	} else if islandlist[4] < 5 {
		buildcpt, err := SetbuildPkt(method, url, 4, islandlist[4]+1, false)
		if err != nil {
			return nil, err
		}
		reqbytes = buildcpt
	}
	return reqbytes, nil
}
