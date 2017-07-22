package mystruct

//Respfromsrv ...
//用于接收响应包
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
