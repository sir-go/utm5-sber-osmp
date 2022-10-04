package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Run() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// set X-Request-Id in nginx.conf.d/conf
	//  location / {
	//      proxy_pass http://upstream;
	//      proxy_set_header X-Request-Id $pid-$msec-$remote_addr-$request_length;
	//  }

	r.Use(
		func(c *gin.Context) { c.Request.URL.RawQuery = strings.ToLower(c.Request.URL.RawQuery) },
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

	r.GET("/", gin.BasicAuth(CFG.Service.Users), SRV.Handler)

	addr := fmt.Sprintf("%s:%d", CFG.Service.Host, CFG.Service.Port)
	LOG.Println("run web-server on http://" + addr)

	s.srv = &http.Server{
		Addr:           addr,
		WriteTimeout:   CFG.Service.Timeouts.Write.Duration,
		ReadTimeout:    CFG.Service.Timeouts.Read.Duration,
		IdleTimeout:    CFG.Service.Timeouts.Idle.Duration,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		eh(err)
	}
}

func (s *Server) Shutdown() {
	LOG.Println("shutdown web-server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	eh(s.srv.Shutdown(ctx))
}
