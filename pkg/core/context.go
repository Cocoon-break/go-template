package core

import (
	stdctx "context"
	"sync"

	"go-template/pkg/trace"
)

type Trace = trace.T

const (
	_TraceName = "_trace_"
)

type StdContext struct {
	stdctx.Context
	Trace
}

var contextPool = &sync.Pool{
	New: func() any {
		return new(context)
	},
}

func NewContext() Context {
	context := contextPool.Get().(*context)
	context.reset()
	return context
}

func ReleaseContext(ctx Context) {
	c := ctx.(*context)
	c.Keys = nil
	contextPool.Put(c)
}

// note: 不是协程安全
type Context interface {
	// Trace 获取 Trace 对象
	Trace() Trace
	SetTrace(trace Trace)
	DisableTrace()
	// RequestContext 获取请求的 context (当 client 关闭后，会自动 canceled)
	RequestContext() StdContext
}

type context struct {
	Keys map[string]any
	mu   sync.RWMutex
}

// RequestContext (包装 Trace + Logger) 获取请求的 context (当client关闭后，会自动canceled)
func (c *context) RequestContext() StdContext {
	return StdContext{
		stdctx.Background(),
		c.Trace(),
	}
}

func (c *context) Trace() Trace {
	t, ok := c.Get(_TraceName)
	if !ok || t == nil {
		return nil
	}

	return t.(Trace)
}

func (c *context) SetTrace(trace Trace) {
	c.Set(_TraceName, trace)
}

func (c *context) DisableTrace() {
	c.SetTrace(nil)
}

func (c *context) Set(key string, value any) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

func (c *context) reset() {
	c.Keys = make(map[string]any, 2)
}
