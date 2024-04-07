package trace

// Response 响应信息
type Response struct {
	Header          interface{} `json:"header,omitempty"`            // Header 信息
	Body            interface{} `json:"body,omitempty"`              // Body 信息
	BusinessCode    int         `json:"business_code,omitempty"`     // 业务码
	BusinessCodeMsg string      `json:"business_code_msg,omitempty"` // 提示信息
	HttpCode        int         `json:"http_code,omitempty"`         // HTTP 状态码
	HttpCodeMsg     string      `json:"http_code_msg,omitempty"`     // HTTP 状态码信息
	CostSeconds     float64     `json:"cost_seconds,omitempty"`      // 执行时间(单位秒)
}
