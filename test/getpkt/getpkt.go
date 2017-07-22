package getpkt

import (
	"fmt"

	"char/markdown/test/mystruct"

	cpt "aladinfun.com/TripleDream/TripleDreamServer/common/libs/crypto"
	tdproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_TripleDream_proto"
	afproto "aladinfun.com/TripleDream/TripleDreamServer/proto/autogen/aladinfun_proto"
)

// ByteToPkt ...
//还原proto结构的数据
func ByteToPkt(reqtype string, rspbytes []byte) (mystruct.Respfromsrv, error) {

	if rspbytes == nil {
		return nil, fmt.Errorf("the rspbyte does not exist!")
	}

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
