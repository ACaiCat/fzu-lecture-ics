package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.New(server.WithHostPorts(":15451"))
	h.Use(Logger())
	h.GET("/v1/lecture/calendar", GetLectureIcs)
	h.Spin()
}
