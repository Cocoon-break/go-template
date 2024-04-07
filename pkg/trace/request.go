package trace

// Request 请求信息
type Request struct {
	TTL        string      `json:"ttl,omitempty"`         // 请求超时时间
	Method     string      `json:"method,omitempty"`      // 请求方式
	DecodedURL string      `json:"decoded_url,omitempty"` // 请求地址
	Header     interface{} `json:"header,omitempty"`      // 请求 Header 信息
	Body       interface{} `json:"body,omitempty"`        // 请求 Body 信息
}
