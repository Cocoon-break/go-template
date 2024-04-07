package rest

import (
	"fmt"
	"net/http"

	"go-template/config"
	"go-template/internal/rest/resp"
	"go-template/pkg/env"
	"go-template/pkg/zlog"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func Start() {
	go startHttpServer()
}

func startHttpServer() {
	basic := config.GetRest()
	gin.SetMode(gin.ReleaseMode)
	// 按照顺序添加插件,其中ResponseBodyBuffer 一定要在ApiLogger之前
	middlewares := []gin.HandlerFunc{
		ginzap.CustomRecoveryWithZap(zlog.GetZapLogger(), true, HandleRecovery),
		gzip.Gzip(gzip.DefaultCompression),
		BizContext(),
	}

	r := gin.Default()
	r.Use(middlewares...)
	r.GET("/version", func(ctx *gin.Context) {
		env := env.CompileInfo()
		resp.JSON(ctx, env)
	})
	r.PUT("/log", func(ctx *gin.Context) {
		l := ctx.Query("level")
		zlog.ChangeLogLevel(l)
		resp.OK(ctx)
	})

	profileGroup := r.Group("/go-template", func(ctx *gin.Context) {
		if ctx.Request.Header.Get("auth") != basic.PprofToken {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Next()
	})
	pprof.RouteRegister(profileGroup, "pprof")
	addr := fmt.Sprintf(":%d", basic.Port)
	zlog.Info("rest", zlog.String("u_msg", fmt.Sprintf("http listen on %s", addr)))
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
