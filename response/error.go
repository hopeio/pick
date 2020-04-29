package response

type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

type ErrCode uint32

const (
	OK  ErrCode = 0
	Err ErrCode = 1
)

var codeStr = map[ErrCode]string{
	OK:  "成功",
	Err: "失败",
}

func (e ErrCode) Error() string {
	if msg, ok := codeStr[e]; ok {
		return msg
	}
	return "未知错误"
}
