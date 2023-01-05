package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

func Run(conf Config) (shutdown func() error) {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	tzLocation, err := time.LoadLocation(conf.Service.Location)
	if err != nil {
		panic(err)
	}

	r.Use(
		func(c *gin.Context) { c.Request.URL.RawQuery = strings.ToLower(c.Request.URL.RawQuery) },
		func(c *gin.Context) { c.Set("config", conf) },
		func(c *gin.Context) { c.Set("tzLocation", tzLocation) },
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s [%s][@%s] %s%3d%s %s%-7s%s %s [%d %s]\n%s",
				param.TimeStamp.Format("2006/01/02 15:04:05.000"),
				param.Request.Header.Get("X-Request-Id"),
				param.Keys["user"],
				param.StatusCodeColor(), param.StatusCode, param.ResetColor(),
				param.MethodColor(), param.Method, param.ResetColor(),
				param.Path,
				param.BodySize,
				param.Latency,
				param.ErrorMessage,
			)
		}),
		gin.Recovery(),
		cors.Default(),
	)

	r.GET("/", gin.BasicAuth(conf.Service.Users), handle)

	addr := fmt.Sprintf("%s:%d", conf.Service.Host, conf.Service.Port)
	//goland:noinspection HttpUrlsUsage
	zlog.Info().Str("address", "http://"+addr).Msg("run web-server")

	srv := &http.Server{
		Addr:           addr,
		WriteTimeout:   conf.Service.Timeouts.Write,
		ReadTimeout:    conf.Service.Timeouts.Read,
		IdleTimeout:    conf.Service.Timeouts.Idle,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zlog.Err(err).Msg("starting web-server")
	}

	return func() error {
		zlog.Info().Msg("shutdown web-server")
		timedCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		return srv.Shutdown(timedCtx)
	}
}
