package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func Logger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 请求开始时间
		start := time.Now()

		// 获取请求信息
		method := string(ctx.Method())
		path := string(ctx.Path())

		ctx.Next(c)
		// 请求结束时间
		end := time.Now()
		latency := end.Sub(start)
		statusCode := ctx.Response.StatusCode()

		// 打印响应日志
		fmt.Printf("[%s] %s %s - %d in %v\n",
			end.Format("2006-01-02 15:04:05"),
			method, path, statusCode, latency)
	}
}
