package setpkt

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"

	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"
)

func SetloginPkt() (io.Reader, error) {
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

	loginreader := bytes.NewReader(loginbytescpt)

	return loginreader, nil
}

func SetbuildPkt() (io.Reader, error) {
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

//SethelloPkt ...
//hello请求
func SethelloPkt() (io.Reader, error) {
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
