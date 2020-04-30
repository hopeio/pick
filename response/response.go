package response

import "net/http"

type ResData struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	//验证码
	Details interface{} `json:"details"`
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
