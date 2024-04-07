package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-template/pkg/core"
	"go-template/pkg/trace"
	"go-template/pkg/zlog"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const GinBizCtxKey = "biz_ctx"

func HandleRecovery(c *gin.Context, err interface{}) {
	path := c.FullPath()
	c.AbortWithStatus(http.StatusInternalServerError)
	zlog.Error("gin_panic", zlog.String("uri", path), zlog.Any("err", err))
}

// RequestId init hte request id
func BizContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// prometheus 监控数据获取不进入
		if strings.Contains(c.Request.URL.RequestURI(), "metrics") {
			c.Next()
			return
		}
		context := core.NewContext()
		defer core.ReleaseContext(context)

		traceId := c.Request.Header.Get("CtxTraceId")

		if traceId == "" {
			traceId = uuid.NewV4().String()
		}
		t := trace.New(traceId)
		decodedURL, _ := url.QueryUnescape(c.Request.URL.RequestURI())
		headers := []string{"timestamp", "authorization"}
		traceHeader := make(map[string]string, len(headers))
		for _, headerKey := range headers {
			if c.GetHeader(headerKey) == "" {
				continue
			}
			traceHeader[headerKey] = c.GetHeader(headerKey)
		}
		t.WithRequest(&trace.Request{
			Method:     c.Request.Method,
			DecodedURL: decodedURL,
			Header:     traceHeader,
			Body:       getRequestBody(c),
		})

		context.SetTrace(t)

		c.Set(GinBizCtxKey, context)
		start := time.Now()
		c.Next()

		t.Success = !c.IsAborted() && (c.Writer.Status() == http.StatusOK)
		t.CostMilliseconds = time.Since(start).Milliseconds()

		zlog.Info("biz_ctx", zlog.Any("trace", t))
	}
}

// 获取请求体
func getRequestBody(ctx *gin.Context) interface{} {
	switch ctx.Request.Method {
	case http.MethodGet:
		return ctx.Request.URL.Query()
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return string(bodyBytes)
	}
	return nil
}
