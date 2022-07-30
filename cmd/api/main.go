// Code generated by hertz generator.

package main

import (
	"github.com/DouYin/cmd/api/rpc"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default(server.WithHostPorts(":8888"))
	rpc.InitRPC()
	register(h)
	h.Spin()
}
