// Package handler defines HTTP handlers and routing configuration,
// including metrics and pprof endpoints.
package handler

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// NewHandler initializes HTTP routes including metrics and pprof endpoints.
func NewHandler() http.Handler {

	handler := gin.New()

	handler.Use(gin.Recovery())

	handler.GET("/metrics", Metrics)

	pprofGroup := handler.Group("/debug/pprof")
	pprofGroup.GET("/", gin.WrapH(http.HandlerFunc(pprof.Index)))
	pprofGroup.GET("/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
	pprofGroup.GET("/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
	pprofGroup.GET("/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
	pprofGroup.GET("/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))

	pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
	pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))

	return handler

}
