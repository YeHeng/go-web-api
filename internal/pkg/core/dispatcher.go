package core

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	"github.com/YeHeng/go-web-api/internal/code"
	"github.com/YeHeng/go-web-api/internal/middleware"
	"github.com/YeHeng/go-web-api/internal/pkg/context"
	"github.com/YeHeng/go-web-api/internal/pkg/logger"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/YeHeng/go-web-api/pkg/errno"
	"github.com/YeHeng/go-web-api/pkg/trace"

	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func dispatcher(r *gin.Engine) {

	withoutTracePaths := map[string]bool{
		"/metrics": true,

		"/debug/pprof/":             true,
		"/debug/pprof/cmdline":      true,
		"/debug/pprof/profile":      true,
		"/debug/pprof/symbol":       true,
		"/debug/pprof/trace":        true,
		"/debug/pprof/allocs":       true,
		"/debug/pprof/block":        true,
		"/debug/pprof/goroutine":    true,
		"/debug/pprof/heap":         true,
		"/debug/pprof/mutex":        true,
		"/debug/pprof/threadcreate": true,

		"/favicon.ico": true,

		"/system/health": true,
	}

	log := logger.Get()
	cfg := config.Get()
	feature := cfg.Feature
	r.Use(func(c *gin.Context) {
		ts := time.Now()

		ctx := context.NewContext(c)
		defer context.ReleaseContext(ctx)

		ctx.Init()

		if !withoutTracePaths[c.Request.URL.Path] {
			if traceId := ctx.GetHeader(trace.Header); traceId != "" {
				ctx.SetTrace(trace.New(traceId))
			} else {
				ctx.SetTrace(trace.New(""))
			}
		}

		defer func() {
			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				log.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", stackInfo))
				ctx.AbortWithError(errno.NewError(
					http.StatusInternalServerError,
					code.ServerError,
					code.Text(code.ServerError)),
				)

				if feature.PanicNotify {
					// notify(context, err, stackInfo)
				}
			}

			if c.Writer.Status() == http.StatusNotFound {
				return
			}

			var (
				response        interface{}
				businessCode    int
				businessCodeMsg string
				abortErr        error
				traceId         string
				graphResponse   interface{}
			)

			if c.IsAborted() {
				for i := range c.Errors { // gin error
					multierr.AppendInto(&abortErr, c.Errors[i])
				}

				if err := ctx.AbortError(); err != nil { // customer err
					multierr.AppendInto(&abortErr, err.GetErr())
					response = err
					businessCode = err.GetBusinessCode()
					businessCodeMsg = err.GetMsg()

					if x := ctx.Trace(); x != nil {
						ctx.SetHeader(trace.Header, x.ID())
						traceId = x.ID()
					}

					c.JSON(err.GetHttpCode(), &code.Failure{
						Code:    businessCode,
						Message: businessCodeMsg,
					})
				}
			} else {
				response = ctx.GetPayload()
				if response != nil {
					if x := ctx.Trace(); x != nil {
						ctx.SetHeader(trace.Header, x.ID())
						traceId = x.ID()
					}
					c.JSON(http.StatusOK, response)
				}
			}

			graphResponse = ctx.GetPayload()

			if feature.RecordMetrics {
				uri := ctx.URI()
				if alias := ctx.Alias(); alias != "" {
					uri = alias
				}

				middleware.RecordMetrics(
					ctx.Method(),
					uri,
					!c.IsAborted() && c.Writer.Status() == http.StatusOK,
					c.Writer.Status(),
					businessCode,
					time.Since(ts).Seconds(),
					traceId,
				)
			}

			var t *trace.Trace
			if x := ctx.Trace(); x != nil {
				t = x.(*trace.Trace)
			} else {
				return
			}

			decodedURL, _ := url.QueryUnescape(c.Request.URL.RequestURI())

			// c.Request.Header，精简 Header 参数
			traceHeader := map[string]string{
				"Content-Type":             c.GetHeader("Content-Type"),
				config.HeaderLoginToken:    c.GetHeader(config.HeaderLoginToken),
				config.HeaderSignToken:     c.GetHeader(config.HeaderSignToken),
				config.HeaderSignTokenDate: c.GetHeader(config.HeaderSignTokenDate),
			}

			t.WithRequest(&trace.Request{
				TTL:        "un-limit",
				Method:     c.Request.Method,
				DecodedURL: decodedURL,
				Header:     traceHeader,
				Body:       string(ctx.RawData()),
			})

			var responseBody interface{}

			if response != nil {
				responseBody = response
			}

			if graphResponse != nil {
				responseBody = graphResponse
			}

			t.WithResponse(&trace.Response{
				Header:          c.Writer.Header(),
				HttpCode:        c.Writer.Status(),
				HttpCodeMsg:     http.StatusText(c.Writer.Status()),
				BusinessCode:    businessCode,
				BusinessCodeMsg: businessCodeMsg,
				Body:            responseBody,
				CostSeconds:     time.Since(ts).Seconds(),
			})

			t.Success = !c.IsAborted() && c.Writer.Status() == http.StatusOK
			t.CostSeconds = time.Since(ts).Seconds()

			log.Info("core-interceptor",
				zap.Any("method", c.Request.Method),
				zap.Any("path", decodedURL),
				zap.Any("http_code", c.Writer.Status()),
				zap.Any("business_code", businessCode),
				zap.Any("success", t.Success),
				zap.Any("cost_seconds", t.CostSeconds),
				zap.Any("trace_id", t.Identifier),
				zap.Error(abortErr),
			)
		}()

		c.Next()
	})
}
