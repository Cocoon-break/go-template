package trace

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
)

var _ T = (*Trace)(nil)

type T interface {
	i()
	ID() string
	WithRequest(req *Request) *Trace
	WithResponse(resp *Response) *Trace
	// 主动发起到外部的请求参数
	WithSendRequest(req *Request) *Trace
	// 主动发送到外部请求到响应结构
	WithSendResponse(resp *Response) *Trace
	AppendSQL(sql *SQL) *Trace
}

// Trace 记录的参数
type Trace struct {
	mux              sync.Mutex
	Identifier       string    `json:"trace_id"`           // 链路ID
	Request          *Request  `json:"request,omitempty"`  // 请求信息
	Response         *Response `json:"response,omitempty"` // 返回信息
	SendRequest      *Request  `json:"send_request,omitempty"`
	SendResponse     *Response `json:"send_response,omitempty"`
	SQLs             []*SQL    `json:"sqls,omitempty"`              // 执行的 SQL 信息
	Success          bool      `json:"success,omitempty"`           // 请求结果 true or false
	CostMilliseconds int64     `json:"cost_milliseconds,omitempty"` // 执行时长(单位豪秒)
}

func New(id string) *Trace {
	if id == "" {
		buf := make([]byte, 10)
		_, _ = io.ReadFull(rand.Reader, buf)
		id = hex.EncodeToString(buf)
	}

	return &Trace{
		Identifier: id,
	}
}

func (t *Trace) i() {}

// ID 唯一标识符
func (t *Trace) ID() string {
	return t.Identifier
}

// WithRequest 设置request
func (t *Trace) WithRequest(req *Request) *Trace {
	t.Request = req
	return t
}

// WithSendRequest 设置request
func (t *Trace) WithSendRequest(req *Request) *Trace {
	t.SendRequest = req
	return t
}

// WithResponse 设置response
func (t *Trace) WithResponse(resp *Response) *Trace {
	t.Response = resp
	return t
}

// WithSendResponse 设置response
func (t *Trace) WithSendResponse(resp *Response) *Trace {
	t.SendResponse = resp
	return t
}

// AppendSQL 追加 SQL
func (t *Trace) AppendSQL(sql *SQL) *Trace {
	if sql == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.SQLs = append(t.SQLs, sql)
	return t
}
