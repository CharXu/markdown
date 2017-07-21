package setpkt

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"
)

//SetloginPkt ...
//生成加密，编码之后的http请求主体
func SetloginPkt() ([]byte, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := strconv.Itoa(r.Intn(100000000))
	loginreq := &tdproto.LoginReq{
		Type:        tdproto.LOGIN_TYPE_GUEST,
		AccessToken: token,
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
	return loginbytescpt, nil
}

func SetbuildPkt() ([]byte, error) {
	buildreq := &tdproto.BuildReq{
		Index:   3,
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
	return buildbytescpt, nil
}

//SethelloPkt ...
//hello请求
func SethelloPkt() ([]byte, error) {
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
	return hellobytescpt, nil
}
