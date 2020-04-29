package response

import "net/http"

type ResData struct {
	Code    uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	//验证码
	Details interface{} `protobuf:"bytes,3,opt,name=details,proto3" json:"details,omitempty"`
}

type File struct {
	File http.File
	Name string
}

type HttpBody struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}
