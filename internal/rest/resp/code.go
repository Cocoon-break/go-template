package resp

import (
	"net/http"

	"go-template/pkg/core"
	"go-template/pkg/trace"

	"github.com/gin-gonic/gin"
)

// 服务正常
func OK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, NewResp(SuccessCode, "ok"))
}

func OkWithBizCtx(ctx *gin.Context, bizCtx core.Context) {
	payload := NewResp(SuccessCode, "ok")
	jsonWithBizCtx(ctx, bizCtx, http.StatusOK, payload)
}

// 服务正常 并返回 JSON 数据结构
func JSON(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, NewResp(SuccessCode, "ok").WithData(v))
}

func JSONWithBizCtx(ctx *gin.Context, bizCtx core.Context, v interface{}) {
	payload := NewResp(SuccessCode, "ok").WithData(v)
	jsonWithBizCtx(ctx, bizCtx, http.StatusOK, payload)
}

func jsonWithBizCtx(ctx *gin.Context, bizCtx core.Context, statusCode int, payload Resp) {
	bizCtx.Trace().WithResponse(&trace.Response{
		Body:     payload,
		HttpCode: statusCode,
	})
	ctx.JSON(statusCode, payload)
}

func ErrWithBizCtx(ctx *gin.Context, bizCtx core.Context, err error) {
	payload := NewResp(FailedCode, err.Error())
	jsonWithBizCtx(ctx, bizCtx, http.StatusOK, payload)
}

// 业务码
func ErrWithBizCtxShouldNotUseCache(ctx *gin.Context, bizCtx core.Context, err error) {
	payload := NewResp(ShouldNotUseCache, err.Error())
	jsonWithBizCtx(ctx, bizCtx, http.StatusOK, payload)
}

// 响应状态码200，同时提供错误信息
func ERRJSON(ctx *gin.Context, err error) {
	ERR(ctx, http.StatusOK, FailedCode, err)
}

// 服务异常
func ERR(ctx *gin.Context, httpCode, code int, arg interface{}) {
	var msg string
	switch arg := arg.(type) {
	case error:
		msg = arg.Error() // arg.(error).Error()
	case string:
		msg = arg // arg.(string)
	default:
		msg = ""
	}
	ctx.JSON(httpCode, NewResp(code, msg))
}

// 请求8参数非法时
func ERRBadReq(ctx *gin.Context, err error) bool {
	if err != nil {
		ERR(ctx, http.StatusBadRequest, FailedCode, err.Error())
		return true
	}
	return false
}

func ERRBadReqWithBizCtx(ctx *gin.Context, bizCtx core.Context, err error) {
	payload := NewResp(FailedCode, err.Error())
	jsonWithBizCtx(ctx, bizCtx, http.StatusBadRequest, payload)
}

func ERRInternalServerBizCtx(ctx *gin.Context, bizCtx core.Context, err error) {
	payload := NewResp(FailedCode, err.Error())
	jsonWithBizCtx(ctx, bizCtx, http.StatusInternalServerError, payload)
}

// 鉴权失败时，响应401
func ERRUnauthorized(ctx *gin.Context, err error) {
	ERR(ctx, http.StatusUnauthorized, UnauthorizedCode, err)
	ctx.Abort()
}
