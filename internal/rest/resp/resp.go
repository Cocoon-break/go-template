package resp

type Resp interface {
	i()                             // i 为了避免被其他包实现
	WithData(data interface{}) Resp // WithData 设置成功时返回的数据
	WithMessage(id string) Resp     // WithMessage 设置当前请求的错误信息
}

type resp struct {
	Code    int         `json:"code"`           // 业务编码
	Message string      `json:"msg,omitempty"`  // 错误描述
	Data    interface{} `json:"data,omitempty"` // 成功时返回的数据
}

const (
	FailedCode        = -1 // 通用失败业务码
	SuccessCode       = 10000
	UnauthorizedCode  = -2 // 未授权
	ShouldNotUseCache = -3 // 不是用缓存
)

func NewResp(code int, msg string) Resp {
	return &resp{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

func (r *resp) i() {}

func (r *resp) WithData(data interface{}) Resp {
	r.Data = data
	return r
}

func (r *resp) WithMessage(msg string) Resp {
	r.Message = msg
	return r
}
