package api

import (
	"github.com/fs714/github-star-manager/api/middleware"
	"github.com/fs714/github-star-manager/api/v1/public"
	"github.com/fs714/github-star-manager/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.Config.Common.RunMode)
	gin.DisableConsoleColor()
	r := gin.New()
	r.Use(middleware.LogWithSkipPath([]string{}))
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	if config.Config.Common.Profiling {
		pprof.Register(r)
	}

	// no authentication
	v1PublicGroup := r.Group("")
	{
		public.InitRoute(v1PublicGroup)
	}

	return r
}
