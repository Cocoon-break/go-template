package req

import (
	"time"

	"go-template/internal/rest"
	"go-template/internal/rest/resp"
	"go-template/pkg/core"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/pkg/errors"
)

func GetBizCtx(ctx *gin.Context) (core.Context, error) {
	var pCtx core.Context
	bizCtx, ok := ctx.Get(rest.GinBizCtxKey)
	if !ok {
		return nil, errors.New("not set biz context")
	}
	pCtx, ok = bizCtx.(core.Context)
	if !ok {
		return nil, errors.New("invalid biz context")
	}
	return pCtx, nil
}

func RateLimit(fillInterval time.Duration, capacity, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, capacity, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			resp.ERRUnauthorized(c, errors.New("rate limit"))
			c.Abort()
			return
		}
		c.Next()
	}
}
